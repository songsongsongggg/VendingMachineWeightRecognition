# VendingMachineWeightRecognition


### pkg/exception/exception.go
定义系统异常类型：
- SensorError: 传感器异常
- ForeignObjectError: 异物异常
- RecognitionError: 识别异常

### pkg/model/model.go
定义基础数据模型：
- Goods: 商品信息
- Stock: 库存信息
- Layer: 层信息

### pkg/recognition/result.go
定义识别结果相关结构：
- RecognitionItem: 识别到的商品
- RecognitionException: 识别异常
- RecognitionResult: 识别结果

### pkg/recognition/weight.go
实现重量识别器：
- WeightRecognizer: 重量识别器结构体
- NewWeightRecognizer: 创建识别器
- Recognize: 识别方法
- recognizeLayer: 单层识别方法

### pkg/recognition/weight_test.go
包含所有测试用例：
- 基础功能测试
- 异常处理测试
- 多商品识别测试
- 边界条件测试

### go.mod
项目依赖管理文件