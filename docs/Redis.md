# Redis数据接口

---
## 说明
```
db: check config/redis
```
---

## 用户
#### userInfo
##### key
|||
|:---|:---|
|type|`hashMap`|
|key|`betahouse:heatae:user:userInfo`|  
|hash key|`userId`|

##### value
一个json string  
example:
```json
{"userInfoId":"201811302141082469393700073084","userId":"201811302141081651290001201884","stuId":"15901116","realName":"王长饶","sex":"男","major":"机械设计制造及其自动化","classId":"15090111","grade":"2015","enrollDate":1441036800000,"extInfo":{}}
```

---

#### roleInfo
##### key

|||
|:---|:---|
|type|`hashMap`|
|key|`betahouse:heatae:user:roleInfo`|  
|hash key|`userId`|

##### value
一个json string
example:
```json

```