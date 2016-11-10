<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

测试客户端
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-09

目录
--------------
###目的
为测试管理员和用户的基本功能提供便利

###请求

* 请求方法: 
* 请求URL: 
* 请求头: 
* 请求参数: 

请求方法、URL、头是固定的, 如需携带的参数均已经预设

###页面设计 
无


###程序实现逻辑
```Title:测试客户端
TestClient-->>YCE: 发送某个用例的各个请求
TestClient<<--YCE: 接受某个用例的响应并判断是否是正确结果
```
大多数请求可以以返回代码作为判断标准, 部分请求需要跟特定的响应数据进行比对是否正确。

###响应数据结构: 
特定请求的响应数据已经预设


### 备注

[ ] 进行单元测试, 确定每个单元都能在正确输入下都能得到正确输出
[ ] 进行基本功能测试, 以能正确执行所有正确输入为目标, 中间不出错, 得到预期结果。
[ ] 进行全路径覆盖测试, 以测试系统的鲁棒性


### 以下为旧版本, 无效///////////////////////////////////////////////////

测试包
-----------

测试包为各个关键功能函数提供测试帮助。包括:

* 一个测试用的rest client
* 一组测试用的json文件 或 常量
* 一组测试相关的方法

其基本数据类型为:

type testclient struct {
    // rest client, for calling API or invoking the method.
    // 调用myhttpclient包
    cli *client
    
    // request describe request
    request 
    
    // response store response for advanced validation
    response 
}

func (*testclient) Validate() bool {
    expect := readFromJson()
    actual := response.Body()
    
    if expect == actual {
        return true 
    } else {
        return false 
    }
}


### 功能性测试
仅验证是否能完成基本功能

### 单元测试
验证多种错误情况、边界条件等, 通过多个测试用例, 来精确函数的行为

### 计划
目前的测试应面向所有的handler展开,输入输出均以handler为准, 辅助函数的准确性暂不讨论

例如要测试功能发布应用, 需要准备的有:

* 待输入的json常量 
* 准备调用该API的client
* 创建后再次获取,将其结果解码为struct1, 与正确json结构的struct2进行比较, 如果温和说明成功

多用例单元测试的时候, 需要先弄清楚多种出错情况, 并针对这些情况设计json常量。其他同上
