(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-2d0c8e41"],{"578a":function(t,s,a){"use strict";a.r(s);var e=function(){var t=this,s=t.$createElement,a=t._self._c||s;return a("div",{staticClass:"login container"},[a("form",{on:{submit:function(s){return s.preventDefault(),t.postData.apply(null,arguments)}}},[t._m(0),a("div",{staticClass:"form-group"},[a("label",[t._v("账号")]),a("input",{directives:[{name:"model",rawName:"v-model",value:t.username,expression:"username"}],staticClass:"form-control",attrs:{type:"text"},domProps:{value:t.username},on:{input:function(s){s.target.composing||(t.username=s.target.value)}}})]),a("div",{staticClass:"form-group"},[a("label",{attrs:{for:"exampleInputPassword1"}},[t._v("密码")]),a("input",{directives:[{name:"model",rawName:"v-model",value:t.password,expression:"password"}],staticClass:"form-control",attrs:{type:"password"},domProps:{value:t.password},on:{input:function(s){s.target.composing||(t.password=s.target.value)}}}),a("small",{staticClass:"form-text text-muted"},[t._v("请记住您的密码")])]),a("div",{staticClass:"row"},[a("div",{staticClass:"col-4"},[a("router-link",{staticClass:"btn btn-link btn-back",attrs:{to:"/Register"}},[t._v(" 注册新用户 ")])],1),t._m(1)])]),a("div",{staticClass:"login-about"},[t._v("作者 QQ759323585 承接软件定制")])])},n=[function(){var t=this,s=t.$createElement,a=t._self._c||s;return a("h1",{staticClass:"headline"},[t._v("登录 "),a("small",[t._v("0.1.0")])])},function(){var t=this,s=t.$createElement,a=t._self._c||s;return a("div",{staticClass:"col-8"},[a("button",{staticClass:"btn btn-dark btn-block",attrs:{type:"submit"}},[t._v("登录")])])}],o=a("c0d6"),r={name:"Login",data:function(){return{username:o["a"].state.userName,password:""}},methods:{postData:function(){o["a"].login(this.username,this.password).then((function(t){console.log("login",t),t.state||alert(t.message)}))}}},l=r,i=a("2877"),u=Object(i["a"])(l,e,n,!1,null,null,null);s["default"]=u.exports}}]);
//# sourceMappingURL=chunk-2d0c8e41.b704ca4e.js.map