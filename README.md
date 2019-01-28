# generator-gin-auth-swag [![NPM version][npm-image]][npm-url] [![Build Status][travis-image]][travis-url] [![Dependency Status][daviddm-image]][daviddm-url] [![Coverage percentage][coveralls-image]][coveralls-url]
> 

## Installation

First, install [Yeoman](http://yeoman.io) and generator-gin-auth-swag using [npm](https://www.npmjs.com/) (we assume you have pre-installed [node.js](https://nodejs.org/)).

```bash
npm install -g yo
npm install -g generator-gin-auth-swag
```

当前项目没有发布到npm，以上命令可能执行失败，可以git clone在本地，然后使用以下命令
```bash
npm link
```

Then generate your new project:

```bash
yo gin-auth-swag
```

进入项目目录，然后运行以下命令
```bash
swag init
```
如果没有swag命令，请执行以下命令
```bash
go get -v -u github.com/swaggo/swag/cmd/swag
cd $GOPATH/src/github.com/swaggo/swag/cmd/swag
go install
```

尝试启动项目
```bash
go run main.go DEV
```

注意，如果localhost没有安装redis服务，模板将在验证auth2.0的时候失败

## Getting To Know Yeoman

 * Yeoman has a heart of gold.
 * Yeoman is a person with feelings and opinions, but is very easy to work with.
 * Yeoman can be too opinionated at times but is easily convinced not to be.
 * Feel free to [learn more about Yeoman](http://yeoman.io/).

## License

MIT © [sunjiaying]()


[npm-image]: https://badge.fury.io/js/generator-gin-auth-swag.svg
[npm-url]: https://npmjs.org/package/generator-gin-auth-swag
[travis-image]: https://travis-ci.org/sunjiaying/generator-gin-auth-swag.svg?branch=master
[travis-url]: https://travis-ci.org/sunjiaying/generator-gin-auth-swag
[daviddm-image]: https://david-dm.org/sunjiaying/generator-gin-auth-swag.svg?theme=shields.io
[daviddm-url]: https://david-dm.org/sunjiaying/generator-gin-auth-swag
[coveralls-image]: https://coveralls.io/repos/sunjiaying/generator-gin-auth-swag/badge.svg
[coveralls-url]: https://coveralls.io/r/sunjiaying/generator-gin-auth-swag
