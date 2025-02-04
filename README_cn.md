# 项目来源
https://bitbucket.org/rj/goey

## 项目目标
1、逐步实现对CEF的集成

## 说明
1、由于golang的限制并在github上统一进行维护，所以将包名全部改成github.com/chaolihf/goey，
2.1、修改TabElement对应的Windows版本，可以通过获取tabitems等对项目进行操作（存在关联子元素无法显示的缺陷，但不影响不包含子项目的情况）；
2.2、将其中容器对应的元素改为大写对外暴露，如vboxElement改为VBoxElement，这样就可以通过window的方法来获取，参考example\controls\main.go相应代码
3、调整github.com/lxn/win到	github.com/chaolihf/win，用来增加缺少的函数来开始重画tab的
3.1、tab增加withCloseButton，初步实现画布