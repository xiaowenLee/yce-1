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
