容器云
================

##### 开场
---

今天,非常荣幸为各位领导和同事分享我们团队近一年来的容器云建设的成果。

相信在座的同学们大多数都听说过容器技术,或者大多听说过docker,是的,我们非常幸运的见证了一个新的时代——容器时代的来临。

据美国权威机构rightscale的报告,现在美国80%的应用部署在云上,其中包括公有云也包括私有云。国内这个比例还很低。

##### 容器
----

容器最早由Google公司提出,迄今为止,已经在Google公司内稳定运行了十多年。每一时刻,Google有大概20亿以上的容器在运行,支撑着全球的搜索服务,邮箱服务,地图服务等。

我们探索和钻研容器技术的目标,就是建立一个像谷歌一样稳定的基础设施。

回顾基础实施演变的历史,在2010年前,我们可以称之为金属时代,我们通过人工的方式维护者物理机,物理机是应用分配的单位,这样造成了很大人工成本和资源浪费。2010开始,Openstack平台发布,开启了云计算的1.0时代,技术人员通过Openstack以更细的粒度划分计算和存储资源,大幅提高了资源的利用率和人工成本。但是,以资源抽象为核心的云1.0有天然的缺陷,因为我们最终要运维或者说运营的使我们应用,而不是这些资源,所以,云2.0重点解决以应用为中心的问题。2013年docker发布,但是2015年容器技术才被业界广泛的研究,2015年是容器技术的元年,庆幸的是我们也是在2015年上半年的时候, chuck让于涛我们开始调研容器技术,因为接触的早和领导的大力支持,我们掌握的容器技术一直在业界都保持相对领先的地位。

上个月华为的全连接大会上,华为的副董事长提出了两个预测:第一,到2025年,85%以上的应用会部署到云上;第二,云的影响远远超过技术本身,还会影响商业模式和人的思维模式,引发一系列的商业革命。


##### 从容器到容器云
----

在我们调研完Docker以及Docker的生态后,发现仅仅Docker或者说容器技术并不能很好解决我们公司当前的问题,(于是涛哥带领我继续探索,进过多方调研和详细的论证,我们向领导chuck打造一个基于容器的私有云,即我们今天要展示的产品——易宝容器云。),如果以容器云的方式打造一个平台,那么将会大大降低我们的运维成本和大大提高研发的效率。

容器云重新定义了计算的边界,在云主机时代,我们需要事先进行严谨规划来提高资源的利用率,我们千方百计的保证我们系统、应用的可用性,可用性是考核我们技术人员工作的最重要的指标之一。在面向资源的云主机时代,只能纵向扩展,而不是快速的水平扩展,针对云主机的运维,人工运维是比较普遍的手段,这造成很多的人工事故。容器云能够让应用告诉迭代,从容器云开始,我们能过做到基础设施或者系统层以下永远可用,这基于应用多副本、健康检查和自愈、多可用区域划分、多计算集群、多数据中心,每个数据中心独立成活等手段,可以实现应用的永远可用。容器云能够让应用自动的按需扩容,在流量激增的时候,容器云会自动为应用扩容,在流量下降时,容器云会自动的为应用缩容,这种弹性扩容不仅自动发生,而且能够在秒级实现扩容、缩容。容器云的另一目标就是构建无人值守的运维平台,让运维工作像流水一样流淌,不用人为的去干预他。

云服务一般分为四个层次,SaaS服务像Gmail等邮箱服务,百度地图等; SAE, Paypal等开放平台属于Paas,我们做的容器云属于CaaS,阿里云、亚马逊云属于IaaS。有了CaaS,会更便利PaaS和SaaS的构建,还会屏蔽底层IaaS的复杂性。

容器云从基础设施开始会改变整个软件交付的生态,就想第一次工业革命的蒸汽机革命一样,蒸汽机革命大幅提高了社会的劳动生产率,让生产力水平有了极大的提升;容器云将技术人员从繁杂的沟通和手工操作中解放出来,会大幅提高研发、运维等技术人员的生产率。

让基础实施如丝般顺滑是我们容器云的使命,也是我们做容器云的最主要的目的。让基础设施自我管理,对用户(运维、研发)透明,就像我们日常生活中的使用水和电一样,即开即用,而不用关心它们是从哪来的。容器云另一个像水电等生活服务的方面是它低廉的成本,我们在使用水电的时候,我们会想到节约,而不是省钱。

应用的全生命周期管理,是云2.0与云1.0最大的区别,容器云一方面提供了永远可用的基础设施,另一方面提供了应用的全生命周期管理,及应用的发布、升级、滚动升级、灰度发布、回滚、弹性扩容、应用下线等整个生命周期的支持。


运维人员、开发人员可以通过容器云随时、随地实现秒级的上线和回滚,可以根据需要自动弹性扩容、缩容。能够轻松实现自主上线,甚至自动上线。

容器云从理论上讲,它以软件的方式定义了计算资源和存储资源,将来,还有软件定义的网络。以软件的方式抽象基础设施,另外,基础设施也是可编程的。《不可变基础设施》、《基础设施即代码》等著名的理论,是我们建设容器云的理论依据。


涛哥,勇哥,我,立尧还有新加入的小鲜肉们,用了一年的时间,探索容器技术的生态,不断尝试,终于打造了一款属于易宝的容器云——易宝容器云,右边是我们的logo,它上面的舵代表我们驾驶容器云这艘巨轮,乘风破浪,扬帆远航。

容器化的应用在易宝容器云上能够在1~30秒内完成上线、回滚、扩容等操作。

易宝容器云能够为公司节省至少一半的硬件成本,对各个产品线来说,比如一年一台虚拟机的成本大概在5000元左右,如果使用容器云,成本最高不会高于2500元。

现在我们容器云部署在两个数据中心,电信和世纪互联,资源池的容量大概是:8640个核,8832G内存,320T的存储,欢迎大家来试用。

随着云技术越来越普及,未来的应用都应该是CloudNative的,云原生应用,即以部署在云上依托云的各种特性的原生应用。他的好处出了高可用、秒级启动等,还有一个显著特性就是,它能够自己管理自己生命周期,就想易宝容器云的管理平台,就是一个云原生应用,只需部署一次,以后升级回滚都可由自身来完成。


这张图是易宝容器云管理系统的部署图,他有app,缓存和数据库三部分组成。

接下来这张图是描述一个应用如何实现多数据中心多副本,一键支付的应用部署A机房的Namespace3下,有4个实例,同样,在B机房的Namespace3下仍然有4个实例。命名空间的好处是,比如namespace1下面的应用出现问题,不会影响到namespace3下的一键支付服务,严格的隔离机制让事故不会蔓延。

这是容器云管理应用全生命周期的示意图,从发布、升级、回滚、下限以及扩容、缩容等。

这个是应用生命周期演进的示意图,我们最先发布了一个v1.0的应用,它有三个实例,我们可以将他扩容至4个实例,还可以缩容到2个实例,升级到v1.1版本,如果升级的v1.1版本有问题,那么就回滚到v1.0版本,最后是应用下线的情况,即将这个应用在容器云中移除。


这个图展示了一个外网的请求怎么调度到应用实例上,我们看到了,一个应用的不同实例会被调度到不同的物理机上,这样,任何一台物理机宕机,我们的应用和服务都不受影响。


我们有两级的负责均衡器,可以根据实际需要修改负载均衡策略。


下面我以nginx为例,来演示如何通过易宝容器云管理平台来发布一个nginx。


下面我给大家汇报一下易宝容器云的落地情况,在三月底,易宝管家就部署在容器云中了,其中在中秋节的活动中,容器云为活动提供了强有力的支持。还有OP的消息推送系统也在4月份部署在了容器云上,北京和湖南的日志中心已经部署完毕,正在跟现有的日志中心合并,北京的日志中心扛过了大数据量的压力测试。这些应用都在容器云中稳定运行了五六个月,没有出现过任何问题。

柏涛的易宝开放平台已经部署在M6的容器云中在测试,测试通过后就可以在生产环境的容器云中部署。还有易宝云盘和集团技术各子系统也会相继部署在容器云上,这个工作我们争取农历年前完成。

虽然容器云已经走到了一个阶段,但是仍有一些难题需要在座的各位同事一起去解决,最困难的是,云的思想对我们已经有的概念和习惯是个冲击,无论是对运维还是研发,这个转变的过程会比较漫长和艰难。还有,目前没有基于容器云的工作流程,这也需要我们一起来探索。还有一个就是磨合,开发与运维的磨合,Devops文化的目的就是打破研发和运维之间的沟通鸿沟,我们的目的也是让研发和运维无障碍沟融,让我们的代码流水线像行云流水般流畅。

经过大家一起努力,我们相信我们会总结出一整套完整而有效的流程和规范,也力争沉淀出属于易宝的Devops文化。


但愿易宝容器云会给各位、给公司带来全新的体验! 谢谢大家!


