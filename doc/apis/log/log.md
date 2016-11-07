<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

日志
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-07

目录
--------------
###目的
按照不同的等级输出日志, 方便了解程序的运行状态。

###请求

* 请求方法: 
* 请求URL: 
* 请求头: 
* 请求参数: 

###页面设计 
无


###程序实现逻辑
无

###响应数据结构: 


### 备注
日志等级

|  日志等级：  |  说明：|
|:---------:|:--------------:|
|FATAL      |致命错误 
|ERROR      |错误 
|WARN       |警告
|INFO       |信息 
|DEBUG      |调试
|TRACE      |追踪 

格式形如: 2016-08-23 15:47:41.850299 [DEBUG] list.go: getDcHost() line 25:  ListDeploymentController getDcHost: server=172.21.1.11  

### 以下为旧版本, 无效///////////////////////////////////////////////////


日志规范
============

### 日志格式

日志格式:

TIME: [LOG_LEVEL] [WHERE] [WHO] [Context]

例如:

2016-08-23 15:47:41.850299 [DEBUG] list.go: getDcHost() line 25:  ListDeploymentController getDcHost: server=172.21.1.11  

#### 时间格式

日期采用 yyyy-mm-dd hh:mm:ss.ms 的格式, 精确到毫秒


#### 日志等级

|  日志等级：  |  说明：|
|:---------:|:--------------:|
|FATAL      |致命错误 
|ERROR      |错误 
|WARN       |警告
|INFO       |信息 
|DEBUG      |调试
|TRACE      |追踪 

具体来说,它们的含义是:

* 致命错误: 它会导致系统终止运行,通过状态控制器立即可见。内存溢出、数组越界、空指针访问等。

* 错误: 系统运行困难, 会影响到用户, 通过状态控制器立即可见, 多数情况下需要人为干预解决。例如数据库无法读写、数据库连接超时等。

* 警告: 意料之外的情况发生,但系统可以继续运行。通过状态控制器立即可见。例如使用了废弃的API、请求页未找到等。

* 信息: 正常且有意义的事情发生,通过状态控制器立即可见,并言简意赅。例如系统启停、会话周期(登入、登出)和有意义的边界事件。

* 调试: 用于调试的信息,描述工作流、非测试函数的进入和退出等。仅打到日志,发布版不包含。

* 追踪: 非常详细的信息,一般不常用,仅打到日志。例如每次循环的状态、对象的全部层级等。


#### 处所

描述日志记录发生的地点: 文件名、函数名、行号


#### 主体

描述日志记录发生的主体: 例如某个Controller等 


#### 上下文

简要描述日志记录发生时,一些关键变量信息等。


### 定义

```golang

type YceLogger struct {
    Logger *log.Logger
}

func New() *YceLogger { 
    return &YceLogger{Logger: log.New(os.Stderr, "", log.Lshortfile)}
}

func (y *YceLogger) Printf() {
    y.Logger.SetPrefix()
    y.Logger.Printf()
}

func (y *YceLogger) Fatalf() {
    y.Logger.SetPrefix()
    y.Logger.Fatalf()
}


```

### 示例:

```
2016-08-23 17:19:48.026 [DEBUG] yce-log1.go:46: main(): main debug log
2016-08-23 17:19:48.026 [FATAL] yce-log1.go:51: main(): main fatal error

```



### 参考资料

[logging-levels-logback-rule-of-thumb-to-assign-log-levels](http://stackoverflow.com/questions/7839565/logging-levels-logback-rule-of-thumb-to-assign-log-levels)

[Log user guide](https://svn.apache.org/repos/asf/commons/proper/logging/tags/LOGGING_1_0_3/usersguide.html)

[Log4j](http://logging.apache.org/log4j/1.2/apidocs/org/apache/log4j/Level.html)