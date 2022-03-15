package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pupilcp/go-excel-exporter/global"
	"github.com/pupilcp/go-excel-exporter/helper/request"
	"github.com/pupilcp/go-excel-exporter/model/dao"
	"github.com/xuri/excelize/v2"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type downloadService struct {
}

type ApiResponseData struct {
	Code int
	Info string
	Data ApiData `json:"data"`
}
type ApiData struct {
	TotalPage int                      `json:"total_page"`
	List      []map[string]interface{} `json:"list"`
	Title     []string                 `json:"title"`
	Key       []string                 `json:"key"`
}

type DownloadCtx struct {
	Ctx              context.Context
	Cancel           context.CancelFunc
	ErrChan          chan error
	ApiReqResultChan chan ApiResponseData
	IsError          bool
	ErrMsg           string
	Wg               sync.WaitGroup
}

type MergeFileCtx struct {
	Ctx         context.Context
	ErrChan     chan error
	FileUrlChan chan string
}

var DownloadService downloadService

var letter = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ", "BA", "BB", "BC", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BK", "BL", "BM", "BN", "BO", "BP", "BQ", "BR", "BS", "BT", "BU", "BV", "BW", "BX", "BY", "BZ"}

func (s *downloadService) HandleDownloadTask(requestUrl, requestMethod string, requestParams interface{}, fileName string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			global.Logger.Errorf("处理下载任务出错：%+v，参数：%+v", err, requestParams)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	dlCtx := &DownloadCtx{
		Ctx:              ctx,
		Cancel:           cancel,
		ErrChan:          make(chan error),
		ApiReqResultChan: make(chan ApiResponseData, 10),
		IsError:          false,
		ErrMsg:           "",
		Wg:               sync.WaitGroup{},
	}
	rand.Seed(time.Now().Unix())
	var tmpFileName = "tmp_" + strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(int(rand.Int63n(100000)))
	page := 1
	dlCtx.Wg.Add(1)
	data := s.getDownloadData(dlCtx, page, requestUrl, requestParams, tmpFileName)
	if data == nil {
		return "", errors.New("获取首页数据为空")
	}
	resp := data.(ApiResponseData)
	if resp.Data.TotalPage > page {
		page = 2
		for page <= resp.Data.TotalPage {
			dlCtx.Wg.Add(1)
			var params = make(map[string]interface{})
			for k, v := range requestParams.(map[string]interface{}) {
				params[k] = v
			}
			params["page"] = page
			go s.getDownloadData(dlCtx, page, requestUrl, params, tmpFileName)
			page += 1
		}
	}
	dlCtx.Wg.Wait()
	if dlCtx.IsError {
		//下载任务失败
		return "", errors.New(dlCtx.ErrMsg)
	}
	//执行完后合并临时文件
	mfCtx := &MergeFileCtx{
		Ctx:         ctx,
		ErrChan:     make(chan error),
		FileUrlChan: make(chan string),
	}
	go s.mergeDownloadFile(mfCtx, fileName, tmpFileName, resp.Data.TotalPage)

	select {
	case err := <-mfCtx.ErrChan:
		return "", err
	case fileUrl := <-mfCtx.FileUrlChan:
		return fileUrl, nil
	}
}

func (s *downloadService) getDownloadData(dlCtx *DownloadCtx, page int, requestUrl string, requestParams interface{}, tmpFileName string) interface{} {
	defer func() {
		if err := recover(); err != nil {
			global.Logger.Errorf("获取下载的数据出错：%+v，参数：%+v", err, requestParams)
			dlCtx.ErrChan <- err.(error)
			dlCtx.Cancel()
		}
		dlCtx.Wg.Done()
	}()
	go func() {
		data, err := request.Post(requestUrl, requestParams, "application/json", 60)
		if err != nil {
			global.Logger.Errorf("请求获取数据出错，url：%s, params：%s，error：%+v, 请求响应：%s", requestUrl, requestParams, err, data)
			dlCtx.ErrChan <- err
			return
		}
		var resp ApiResponseData
		dataStr := []byte(data)
		err = json.Unmarshal(dataStr, &resp)
		if err != nil {
			global.Logger.Errorf("解析json报错，url：%s, params：%s，data: %s, error：%+v", requestUrl, requestParams, dataStr, err)
			dlCtx.ErrChan <- err
			return
		}
		if len(resp.Data.List) <= 0 {
			errMsg := fmt.Sprintf("获取数据为空，url：%s, params：%s", requestUrl, requestParams)
			global.Logger.Error(errMsg)
			dlCtx.ErrChan <- errors.New(errMsg)
			return
		}
		//写入临时文件
		s.writeDownloadFile(dlCtx, resp, page, tmpFileName)
		dlCtx.ApiReqResultChan <- resp
	}()

	select {
	case <-dlCtx.Ctx.Done():
		// 其他RPC调用调用失败
		dlCtx.IsError = true
		return nil
	case err := <-dlCtx.ErrChan:
		// 本RPC调用失败，返回错误信息
		// 取消其它任务执行
		dlCtx.IsError = true
		dlCtx.ErrMsg = err.Error()
		dlCtx.Cancel()
		return nil
	case resp := <-dlCtx.ApiReqResultChan:
		// 本RPC调用成功，不返回错误信息
		return resp
	}
}

func (s *downloadService) writeDownloadFile(dlCtx *DownloadCtx, resp ApiResponseData, page int, tmpFileName string) {
	defer func() {
		if err := recover(); err != nil {
			global.Logger.Errorf("分页数据写入临时文件出错：%+v", err)
			dlCtx.ErrChan <- err.(error)
		}
	}()
	f := excelize.NewFile()
	for i, t := range resp.Data.Title {
		f.SetCellValue("Sheet1", letter[i]+"1", t)
	}
	num := 2
	for _, item := range resp.Data.List {
		for n, key := range resp.Data.Key {
			f.SetCellValue("Sheet1", letter[n]+strconv.Itoa(num), item[key])
		}
		num++
	}
	tmpFilePath := global.Config.GetString("log.logPath") + string(os.PathSeparator) + time.Now().Format("2006/01/02")
	_, err := os.Stat(tmpFilePath)
	if err != nil {
		os.MkdirAll(tmpFilePath, os.ModePerm)
	}
	if err := f.SaveAs(tmpFilePath + "/" + tmpFileName + "_" + strconv.Itoa(page) + ".xlsx"); err != nil {
		global.Logger.Errorf("写入文件报错：%+v", err)
		dlCtx.ErrChan <- err
	}
}

func (s *downloadService) mergeDownloadFile(mfCtx *MergeFileCtx, fileName, tmpName string, totalPage int) {
	defer func() {
		if err := recover(); err != nil {
			global.Logger.Errorf("合并临时分页文件出错：%+v", err)
			mfCtx.ErrChan <- err.(error)
		}
	}()
	var err error
	tmpFilePath := global.Config.GetString("log.logPath") + string(os.PathSeparator) + time.Now().Format("2006/01/02")
	output := excelize.NewFile()
	rowNum := 1
	for i := 1; i <= totalPage; i++ {
		tmpFile := tmpFilePath + "/" + tmpName + "_" + strconv.Itoa(i) + ".xlsx"
		_, err = os.Stat(tmpFile)
		if err != nil {
			continue
		}
		//读取临时文件的内容
		f, _ := excelize.OpenFile(tmpFile)
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			continue
		}
		for index, row := range rows {
			if i > 1 && index == 0 {
				continue
			}
			for colNum, colCell := range row {
				//写入
				output.SetCellValue("Sheet1", letter[colNum]+strconv.Itoa(rowNum), colCell)
			}
			rowNum++
		}
		//删除临时文件
		os.Remove(tmpFile)
		f.Close()
	}
	subPath := string(os.PathSeparator) + time.Now().Format("2006/01/02")
	filePath := global.Config.GetString("system.downloadPath") + subPath
	_, err = os.Stat(filePath)
	if err != nil {
		os.MkdirAll(filePath, os.ModePerm)
	}
	rand.Seed(time.Now().Unix())
	downloadFileName := fileName + "_" + time.Now().Format("20060102150405") + strconv.FormatInt(rand.Int63n(99999999), 10) + ".xlsx"
	fileFullPath := filePath + "/" + downloadFileName
	if err := output.SaveAs(fileFullPath); err != nil {
		global.Logger.Errorf("合并写入文件报错：%+v", err)
		mfCtx.ErrChan <- err
	}

	mfCtx.FileUrlChan <- global.Config.GetString("system.domain") + "/download" + subPath + "/" + downloadFileName
}

func (s *downloadService) HandleTask() {
	for taskId := range TaskChannel {
		go func(taskId int) {
			taskInfo := dao.DownloadTaskDao.GetOne(taskId)
			if taskInfo.TaskId <= 0 {
				return
			}
			var requestParams = make(map[string]interface{})
			json.Unmarshal([]byte(taskInfo.RequestParams), &requestParams)
			DownloadTaskService.UpdateTask(taskInfo.TaskId, "", "", 1)
			fileUrl, err := s.HandleDownloadTask(taskInfo.RequestUrl, taskInfo.RequestMethod, requestParams, taskInfo.FileName)
			taskStatus := 2
			errMsg := ""
			if err != nil {
				taskStatus = 3
				errMsg = err.Error()
			}
			DownloadTaskService.UpdateTask(taskInfo.TaskId, fileUrl, errMsg, taskStatus)
		}(taskId)
	}
}
