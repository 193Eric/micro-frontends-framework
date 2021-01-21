## 简介

为了解决庞大的一整块后端服务带来的变更与扩展方面的限制，出现了微服务架构（Microservices）：

> 微服务是面向服务架构（SOA）的一种变体，把应用程序设计成一系列松耦合的细粒度服务，并通过轻量级的通信协议组织起来
具体地，将应用构建成一组小型服务。这些服务都能够独立部署、独立扩展，每个服务都具有稳固的模块边界，甚至允许使用不同的编程语言来编写不同服务，也可以由不同的团队来管理  

然而，越来越重的前端工程也面临同样的问题，自然地想到了将微服务思想应用（照搬）到前端，于是有了「微前端（micro-frontends）」的概念：

>即，一种由独立交付的多个前端应用组成整体的架构风格。具体的，将前端应用分解成一些更小、更简单的能够独立开发、测试、部署的小块，而在用户看来仍然是内聚的单个产品：

简单来说就是将一个巨无霸（Monolith）的前端工程拆分成一个一个的小工程 （主要是ToB业务）

## 特点

### 1.简单、松耦合的代码库
比起一整块的前端代码库，微前端架构下的代码库倾向于更小/简单、更容易开发。 如果一个大型的toB应用，里面有几百个页面，几百个接口，上百个组件，编译起来，就是一场灾难。这个时间意味着什么？它不仅会耽误我们开发人员的时间，还会影响整个团队的效率。上线时，在 Docker、CI 等环境下，耗时还会被延长。如果部署后出几个 Bug，要线上立即修复，那就不知道要熬到几点了。  

使用微前端改造后可以做到编译时间很短，且功能相对独立。测试debug的时间也会变得少，更容易定位。


### 2.优雅的处理历史项目
在开发的过程中，总会有一些不是很美好的代码。正常情况我们看到这个，脑子里面就是两个字（”**重构**“）。

- 历史项目，祖传代码
- 交付压力，当时求快
- 就近就熟，当时求稳

而要对这些代码进行彻底重构的话，最大的问题是很难有充裕的资源去大刀阔斧地一步到位，在逐步重构的同时，既要确保中间版本能够平滑过渡，同时还要持续交付新特性。这个时候如果用微前端的方式，把这些项目独立开来，单独运行。然后再找时间，单独重构这个模块。就能平滑过渡了。




### 3.单独部署打包

独立部署的能力在微前端体系中至关重要，能够缩小变更范围，进而降低相关风险

因此，每个微前端都应具备有自己的持续交付流水线（包括构建、测试并部署到生产环境），并且要能独立部署，不必过多考虑其它代码库和交付流水线的当前状态。

运用微前端的方式，我们可以每个独立的模块，单独的git库。这样的话，即使是别的模块出了问题，也不会影响当前发布的模块。


## 如何构建微前端



![在这里插入图片描述](https://img-blog.csdnimg.cn/20210120170924473.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzI0MDczODg1,size_16,color_FFFFFF,t_70)
### 加载器

Systemjs是一个可配置模块加载器，为浏览器和NodeJs启用动态的Es模板加载器。任何具有标准的URL都可被加载为一个模块：

```javascript
<script src="system.js"></script>



<script>



  // 加载相对于当前地址的url

  SystemJS.import('./local-module.js');


  // 加载绝对url的地址

  SystemJS.import('https://code.jquery.com/jquery.js');



</script>
```

可以加载任何类型的模块格式，并由SystemJS自动检测。



### 包装器

`single-spa`是一个用于前端微服务化的`JavaScript`前端解决方案。

特点：

- (兼容各种技术栈)在同一个页面中使用多种技术框架(React, Vue, AngularJS, Angular, Ember等任意技术框架),并且不需要刷新页面.
- (无需重构现有代码)使用新的技术框架编写代码,现有项目中的代码无需重构.
- (更优的性能)每个独立模块的代码可做到按需加载,不浪费额外资源.每个独立模块可独立运行.

`import-map` 的实现是 `single-spa` 的核心部分，我们首先得知晓什么是 `import-map` 。

我们先来看两段代码

```text
import moment from 'moment';
import 'http://momentjs.com/downloads/moment.js';
```

这是一个典型的 `es6` 的 `ES Module`的代码，但是现在浏览器并不支持。为了让浏览器支持，我们需要对其进行兼容处理。

在一个文件中我们写入如上代码，显然第一行是无法正常运行的，第二行是可以正常运行的，但如果我们想要第一行正常运行的话， `import-map` 就可以粉墨登场啦。只需要在 `html` 文件书写如下：

```text
<script type="importmap">
    {
        "imports": {
            "moment": "https://momentjs.com/downloads/moment.js",
        }
    }
</script>
```

就可以了。

但是现在浏览器并不支持，想要让它支持的话，需要引入 `system.js` 。

```text
<script type="systemjs-importmap">
    {
        "imports": {
            "moment": "https://momentjs.com/downloads/moment.js",
        }
    }
</script>
<script src="https://cdn.jsdelivr.net/npm/systemjs/dist/system.js"></script>
```

而在 `single-spa` 的使用过程中，我们需要用 `import-map` 在根项目中引入所有的模块文件和子项目，从而在其余项目中可以进行模块的引用，就像上面说的那样，可以把 `moment` 想象成一个子项目。

```text
<script type="systemjs-importmap">
    {
        "imports": {
            "module": "https://[cdn-link].js",
        }
    }
</script>
<script src="https://cdn.jsdelivr.net/npm/systemjs/dist/system.js"></script>
```



### 消息总线
应用微服务化之后,每一个单独的模块都是一个黑盒子, 里面发生了什么,状态改变了什么,外面的模块是无从得知的. 比如`模块A`想要根据`模块B`的某一个内部状态进行下一步行为的时候,黑盒子之间没有办法通信.这是一个大麻烦.



每一个模块之间都是有生命周期的.当模块被卸载的时候,如何才能保持后续的正常的通信?



icestarck v1.2.0 推出了正式的官方跨应用通信方案 [@ice/stark-data](https://link.zhihu.com/?target=https%3A//www.npmjs.com/package/%40ice/stark-data)。核心只有两个 API：store(全局变量管理中心) 和 event（全局事件管理中心）。我们针对全局变量，增加了对应的监听事件，基本能覆盖 95% 以上的实际通信场景。简单通过几个场景说明具体使用方式。

```js
import { store } from '@ice/stark-data';

// 子应用A设置信息
store.set('from', 'A');
store.set('current', 0);

// 子应用B读取信息
const from  = store.get('from');
const current = store.get('current');
```


- 子应用页面切换参数流转
```
import { store } from '@ice/stark-data';

// 子应用A设置信息
store.set('from', 'A');
store.set('current', 0);

// 子应用B读取信息
const from  = store.get('from');
const current = store.get('current');
```


- 框架应用顶部有 ”消息“ 展示入口，子应用内有阅读消息的能力，阅读完消息后需要通知框架应用刷新“消息”展示信息

```js
// 框架应用
import { event } from '@ice/stark-data';

function fresh(needFresh) {
  if (!needFresh) return;

  fetch('/api/fresh/message').then(res => {
    // ...
  });
}

event.on('freshMessage', fresh);

// 子应用
import { event } from '@ice/stark-data';

event.emit('freshMessage', false);
// ...
event.emit('freshMessage', true);
```


### 构建部署

每个子项目单独打包上传到主项目同级目录。

用户访问主项目index.html后，js加载器会加载apps.config.js。

无论路由是什么，每次必会首先加载主项目，再根据路由来匹配要加载哪个子项目。

在资源服务器上起一个监听服务(我使用的是nodejs脚本+pm2守护)，原有子项目的部署方式完全不变(前后端完全分离，资源带hash)，当监听服务检测到文件改动时，去子项目部署文件夹里找它的app.js，写入apps.config.js。






### 讨论



router和vuex是可以选择不暴露的，尤其是vuex，我觉得增加了太多成本，也大大降低了未来同一架构接入其他技术栈的可能性（比如react）

数据管理还是由子项目自行处理的，耦合到主项目上不便于后续的拓展。vuex每个子项目有它自己的store实例。

如果有react/ng等其他技术栈，需要给他们找合适的包装器（single-spa-react （我用的包装器就是single-spa-vue）），理论上是可行的，整套架构除了包装器以外都不用动。

权限 是子系统自己单独管理的，逻辑放在子系统公共组件包里



### why not iframe

1. url 不同步。浏览器刷新 iframe url 状态丢失、后退前进按钮无法使用。
2. UI 不同步，DOM 结构不共享。想象一下屏幕右下角 1/4 的 iframe 里来一个带遮罩层的弹框，同时我们要求这个弹框要浏览器居中显示，还要浏览器 resize 时自动居中..
3. 全局上下文完全隔离，内存变量不共享。iframe 内外系统的通信、数据同步等需求，主应用的 cookie 要透传到根域名都不同的子应用中实现免登效果。
4. 慢。每次子应用进入都是一次浏览器上下文重建、资源重新加载的过程。



其中有的问题比较好解决(问题1)，有的问题我们可以睁一只眼闭一只眼(问题4)，但有的问题我们则很难解决(问题3)甚至无法解决(问题2)，而这些无法解决的问题恰恰又会给产品带来非常严重的体验问题， 最终导致我们舍弃了 iframe 方案。