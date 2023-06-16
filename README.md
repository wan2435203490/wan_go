# wan_go

集成go项目

jwt加密userId，生成tokenString,前端请求放在header["Token"]里面。

session采用cookie store存储用户信息curUser。如果session里面没有获取到curUser,将解析token，jwt解码获取userId，查询db获取user。

