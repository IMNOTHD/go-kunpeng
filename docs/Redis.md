# Redis数据接口

---
## 说明
```
db: check config/redis
```
---

## user
#### userInfo
##### key

|type|`hashMap`|
|:---|:---|
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

|type|`set`|
|:---|:---|
|key|`betahouse:heatae:user:roleInfo:userId`|  

##### value
多个描述roleInfo的string  
example of one:
```json
"ACTIVITY_STAMPER"
```

---

#### jobInfo
##### key

|type|`set`|
|:---|:---|
|key|`betahouse:heatae:user:jobInfo:userId`|  

##### value
多个描述jobInfo的string  
example of one:
```json
"时雨技术交流与支持协会社长"
```

---

## activity
### activity record
##### key

|type|`zset`|
|:---|:---|
|key|`betahouse:heatae:activity:activityRecord:{userId}:{activityType}`|  

##### value
多个描述活动章的json string  
example of one:
```json
{"activityRecordId":"201903202000435094633110022019","activityId":"201903201351484863392210012019","userId":"201811302142192259540001201847","scannerUserId":"201811302141557664490001201843","time":0,"type":"lectureActivity","status":"ENABLE","term":"2018B","grades":"","extInfo":{"scannerName":"黄奕雯"},"createTime":1553083243000,"activityName":"《驴得水》话剧公演","organizationMessage":"大学生艺术团","location":null,"startTime":1553078700000,"endTime":1553085900000,"score":null,"activityTime":"0.0","scannerName":"黄奕雯"}
```