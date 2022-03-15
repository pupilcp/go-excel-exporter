# go-excel-exporter
go将业务数据导出excel

### 简介
使用go多协程，通过URL回调的方式，异常获取业务数据并将数据写入excel后进行下载，真正与业务解耦。

### 功能实现
1. 使用go多协程通过回调URL按分页获取业务数据，将每页数据写入临时的excel文件，最后通过合并临时文件，生成最终要下载的文件。

### 特点
1. 使用gin框架
2. 使用go多协程并发分页获取业务数据，提高导出效率
3. 通过channel异步处理要导出excel的任务
4. 失败的导出任务可以进行重试

### 核心依赖库
1. viper解析toml配置文件
2. gorm操作mysql
3. logrus+lumberjack实现日志的记录以及文件的自动切割

### 安装前准备
#### 环境要求
1. go版本最好 >=1.15
2. mysql

#### 导入sql数据
将config/data.sql文件导入到mysql数据库

#### 修改配置文件
修改根目录下的config/config.toml配置文件里的参数，具体参数请参考文件内的说明。

### 安装：
1. git clone https://github.com/pupilcp/go-excel-exporter.git
2. 安装依赖：cd $PATH, go mod tidy
3. 编译：cd $PATH, go build -o server . ，生成server二进制执行文件。

### 使用
1. ./server 启动服务

### 可用接口列表

1. 获取下载任务列表： 域名+/api/task/get-list，请求方式：GET
   请求参数：
   
   | 参数 | 类型 | 说明 |
   | ------ | ------ | ------ |
   | user_id | int | 业务系统的用户ID，用于查询该用户所有的下载任务 |
   | page | int | 下载的文件名 |
   | page_size | int | 请求回调的url |
   
   响应示例：

| 参数 | 类型 | 说明 |
| ------ | ------ | ------ |
| code | int | 业务响应码，1是成功，其它为失败 |
| data | object | 数据内容对象 |
| data.list | array | 下载任务列表 |
| data.list.task_id | int | 下载任务ID |
| data.list.user_id | int | 下载用户ID |
| data.list.file_name | string | 下载文件名 |
| data.list.request_url | string | 回调的url |
| data.list.request_params | string | 回调url的参数 |
| data.list.task_status | int | 下载任务状态，0：未处理，1：处理中，2：成功，3：失败 |
| data.list.remark | string | 备注信息，失败原因会记录到这里 |
| data.list.task_status | int | 下载任务ID |
| data.list.created_at | int | 创建任务的时间戳 |
| data.list.updated_at | int | 更新的时间戳 |
| data.total | int | 数据总条数 |
| info | string | 提示信息 |

```json
   {
   "code": 1,
   "data": {
      "list": [
         {
            "task_id": 25,
            "user_id": 123,
            "file_name": "订单导出",
            "request_url": "http://order.internal.homary.com/order-query/export-order-list",
            "request_params":"{\"add_id\":49,\"add_name\":\"江鹏\",\"admin_id\":49,\"admin_name\":\"江鹏\",\"business_type\":\"\",\"created_at_end\":\"\",\"created_at_start\":\"2022-02-25 00:00:00\",\"email\":\"\",\"is_delivery_method\":\"\",\"is_dispatched\":\"\",\"is_export\":1,\"is_oversold\":\"\",\"is_paid\":\"\",\"is_refunded\":\"\",\"language\":\"zh\",\"max_price\":\"\",\"min_price\":\"\",\"operator\":\"{\\\"type\\\":2,\\\"id\\\":49,\\\"name\\\":\\\"\\\江\\\鹏\\\"}\",\"order_no\":\"\",\"order_sales_type\":\"\",\"order_status\":\"\",\"order_type\":\"\",\"page\":1,\"page_size\":100,\"payment_method_code\":\"\",\"payment_no\":\"\",\"pf\":\"export\",\"phone\":\"\",\"platform\":\"\",\"site_code\":\"\",\"sku_code\":\"\",\"spu_title\":\"\",\"support_after_sale\":\"\",\"time_zone_code\":\"UTC+8\",\"token\":\"ZSyaBIhUShWDUgmtyn4jNa8eOinpsGax\",\"user_id\":\"\",\"user_name\":\"\",\"user_type\":\"\"}",
            "request_method": "POST",
            "download_file": "",
            "task_status": 3,
            "remark": "",
            "created_at": 1646840895,
            "updated_at": 1646840895
         }
      ],
      "total": 89
   },
   "message": "success"
}
```

2. 创建下载任务： 域名+/api/task/create，请求方式：POST，Content-Type: application/json
   请求参数：
   
| 参数 | 类型 | 说明 |
| ------ | ------ | ------ |
| user_id | int | 业务系统的用户ID，用于查询该用户所有的下载任务 |
| file_name | string | 下载的文件名 |
| request_url | string | 请求回调的url |
| request_params | string | 请求回调url的参数 |
   
```json
   {
    "user_id": 123,
    "file_name": "订单导出",
    "request_url": "http://order.internal.homary.com/order-query/export-order-list",
    "request_params": {
        "page":1,
        "page_size":100,
        "order_no":"",
        "payment_no":"",
        "created_at_start": "2021-08-25 00:00:00",
        "created_at_end": ""
    }
}
```
   响应示例：
```json
   {
   "code": 1,
   "info": "success"
}
```

重点注意事项：
参数中的request_url，请求方式：POST，Content-Type: application/json，go-excel-exporter会将request_params参数作为body请求request_url，请求响应的数据结构【必须】严格遵从如下格式：

| 参数 | 类型 | 说明 |
| ------ | ------ | ------ |
| code | int | 业务响应码，1是成功，其它为失败 |
| data | object | 数据内容对象 |
| data.title | array | excel每一列的标题 |
| data.key | array | 每一行数据的下标字段 |
| data.list | array | 数据集合 |
| info | string | 提示信息 |
```json
{
    "code": 1,
    "data": {
        "title": [
            "订单ID",
            "订单类型"
        ],
        "key": [
            "order_id",
            "order_type"
        ],
        "total_page": 119,
        "list": [
            {
                "order_id": "23094",
                "order_type": "普通订单"
            }
        ]
    },
    "info": "success"
}
```

3. 失败重试任务：域名+/api/task/retry，请求方式：POST，content-type:form-data
   请求参数：
   | 参数 | 类型 | 说明 |
   | ------ | ------ | ------ |
   | task_id | int | 下载任务id |

   响应示例：
```json
   {
   "code": 1,
   "info": "success"
}
```  

   

### 支持
go

### 其它
如需交流，请邮件联系：310976780@qq.com