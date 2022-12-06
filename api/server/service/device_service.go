package service

import (
	"go-backend/api/server/dao"
	"go-backend/api/server/entity"
	"go-backend/api/server/tools/server"
	"go-backend/api/server/vo"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/spf13/viper"
)

// @Summary API of golang gin backend
// @Tags Device-fixed
// @description create fixed device  : 创建固定式设备 参数列表：[设备所在的牧舍ID、厂家提供的设备编号、设备类型] 访问携带token
// @version 1.0
// @accept mpfd
// @param CompanyId formData string true "company id"
// @param DeviceId formData string true "device id"
// @param TypeId formData string true "type name"
// @param InstallTime formData string true "install time"
// @param BoughtTime formData string true "bought time"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/fixed/create [post]
func CreateFixedDeviceService(deviceId string, farmhouseId uint, fixedDeviceTypeId string, owner uint, installTime time.Time, boughtTime time.Time) uint {
	fixedDevice := entity.FixedDevice{
		DeviceID: deviceId,
		FarmhouseID: farmhouseId,
		FixedDeviceTypeID: fixedDeviceTypeId,
		Owner: owner,
		BoughtTime: boughtTime,
		InstallTime: installTime,
		Stat: "normal",
	}
	id := dao.CreateFixedDevice(fixedDevice)
	return id
}

// @Summary API of golang gin backend
// @Tags Device-fixed
// @description delete fixed device : 删除固定式设备 参数列表：[设备ID] 访问携带token
// @version 1.0
// @accept application/json
// @param Id query int true "id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/fixed/delete [delete]
func DeleteFixedDeviceService(fixedDeviceId uint) {
	dao.DeleteFixedDevice(fixedDeviceId)
}

// @Summary API of golang gin backend
// @Tags Device-portable
// @description create portable device : 创建便携式设备 参数列表：[设备绑定的生物ID、厂家提供的设备编号、设备类型] 访问携带token
// @version 1.0
// @accept mpfd
// @param BiologyId formData string true "biology id"
// @param DeviceId formData string true "device id"
// @param TypeId formData string true "type name"
// @param InstallTime formData string true "install time"
// @param BoughtTime formData string true "bought time"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/portable/create [post]
func CreatePortableDeviceService(biologyId uint, portableDeviceId string, portableDeviceTypeId string, installTime time.Time, boughtTime time.Time) uint {
	biology := dao.GetBiologyInfoById(biologyId)
	if biology.ID == 0 {
		panic("biology id does not exists")
	}
	owner := biology.Owner
	portableDevice := entity.PortableDevice{
		DeviceID: portableDeviceId,
		BiologyID: biologyId,
		PortableDeviceTypeID: portableDeviceTypeId,
		Owner: owner,
		BoughtTime: boughtTime,
		InstallTime: installTime,
		Stat: "normal",
	}
	id := dao.CreatePortableDevice(portableDevice)
	return id
}

// @Summary API of golang gin backend
// @Tags Device-portable
// @description delete portable device : 删除便携式设备 参数列表：[设备ID] 访问携带token
// @version 1.0
// @accept application/json
// @param Id query int true "id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/portable/delete [delete]
func DeletePortableDeviceService(portableDeviceId uint) {
	dao.DeletePortableDevice(portableDeviceId)
}

// @Summary API of golang gin backend
// @Tags Device-fixed
// @description create fixed device type : 新增固定式设备类型 参数列表：[设备类型]
// @version 1.0
// @accept mpfd
// @param FixedDeviceTypeId formData string true "type name"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/fixed/create_type [post]
func CreateFixedDeviceTypeService(fixedDeviceTypeId string) {
	fixedDeviceType := entity.FixedDeviceType{
		FixedDeviceTypeID: fixedDeviceTypeId,
	}
	dao.CreateFixedDeviceType(fixedDeviceType)
}

// @Summary API of golang gin backend
// @Tags Device-fixed
// @description delete fixed device type : 删除固定式设备类型 参数列表：[设备类型ID] 访问携带token
// @version 1.0
// @accept application/json
// @param FixedDeviceTypeId query int true "type name"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/fixed/delete_type [delete]
func DeleteFixedDeviceTypeService(fixedDeviceTypeId string) {
	dao.DeleteFixedDeviceType(fixedDeviceTypeId)
}

// @Summary API of golang gin backend
// @Tags Device-portable
// @description create portable device type : 新增便携式设备类型 参数列表：[设备类型] 
// @version 1.0
// @accept mpfd
// @param PortableDeviceTypeId formData string true "type name"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router  /device/portable/create_type [post]
func CreatePortableDeviceTypeService(portableDeviceTypeId string) {
	portableDeviceType := entity.PortableDeviceType{
		PortableDeviceTypeID: portableDeviceTypeId,
	}
	dao.CreatePortableDeviceType(portableDeviceType)
}

// @Summary API of golang gin backend
// @Tags Device-portable
// @description delete portable device type : 删除便携式设备类型 参数列表：[设备类型ID] 访问携带token
// @version 1.0
// @accept application/json
// @param PortableDeviceTypeId query int true "type name"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/portable/delete_type [delete]
func DeletePortableDeviceTypeService(fixedDeviceTypeId string) {
	dao.DeleteFixedDeviceType(fixedDeviceTypeId)
}

var streamAccessToken = map[string]interface{} {
	"accessToken" : "null",
	"expireTime" : time.Now().AddDate(0, 0, -2),
}

func updateMonitorStreamToken() {
	appKey := viper.GetString("monitor.AppKey")
	secret := viper.GetString("monitor.Secret")
	url := viper.GetString("monitor.AccessMonitorTokenUrl")
	contentType := "application/x-www-form-urlencoded"
	payload := strings.NewReader("appKey=" + appKey + "&appSecret=" + secret)
	response, err := http.Post(url, contentType, payload)
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		panic("getStreamFailed")
	}
	responseBody := server.GetResponseAccessMonitor(response)
	data := responseBody
	streamAccessToken["accessToken"] = data.AccessToken
	streamAccessToken["expireTime"] = time.Now()
}

// @Summary API of golang gin backend
// @Tags Device-fixed
// @description get monitor streaming address : 获取摄像头的直播地址 参数列表：[摄像头设备ID] 访问携带token
// @version 1.0
// @accept application/json
// @param Id query string true "id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/fixed/get_monitor [get]
func GetMonitorStreamByDeviceIdService(deviceId uint) (string, int64, string, string, string){
	accessToken := streamAccessToken["accessToken"].(string)
	expireTime := streamAccessToken["expireTime"].(time.Time)
	if accessToken == "null" || 
		time.Since(expireTime) >= 24 * time.Hour {
			updateMonitorStreamToken()
		}
	accessToken = streamAccessToken["accessToken"].(string)
	monitorDeviceId := dao.GetFixedDeviceInfoById(deviceId).DeviceID
	serverUrl := viper.GetString("monitoR.GetStreamAddressUrl")
	contentType := "application/x-www-form-urlencoded"
	payload := strings.NewReader("accessToken=" + accessToken + "&deviceSerial=" + monitorDeviceId)
	response, err := http.Post(serverUrl, contentType, payload)
	if err != nil {
		panic(err.Error())
	}
	code := response.StatusCode
	defer response.Body.Close()
	if code != 200 {
		panic("error occurs during get streaming address from server by deviceID ")
	}
	responseBody := server.GetResponseBodyMonitor(response)
	msg := responseBody.Msg
	data := responseBody.Data
	resultUrl := data.Url
	resultExpireTime := data.ExpireTime
	id := data.Id
	return resultUrl, resultExpireTime, id, msg, accessToken
}

var newCollarAccessToken = map[string]interface{} {
	"accessToken" : "null",
	"expireTime" : time.Now().AddDate(0, 0, -2),
}

func updateNewCollarToken() {
	userName := viper.GetString("newcollar.uname")
	password := viper.GetString("newcollar.psw")
	url := viper.GetString("newcollar.login-url")
	contentType := "application/x-www-form-urlencoded"
	payload := strings.NewReader("username=" + userName + "&password=" + password)
	response, err := http.Post(url, contentType, payload)
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		panic("getNewCollarAccessTokenFailed")
	}
	token := server.GetResponseNewCollarAccessToken(response)
	newCollarAccessToken["accessToken"] = "Bearer " + token
	newCollarAccessToken["expireTime"] = time.Now()
}

// @Summary API of golang gin backend
// @Tags Device-portable
// @description get new-type collar realtime data by device id : 获取中农智联项圈的最新数据 参数列表：[设备ID] 访问携带token
// @version 1.0
// @accept application/json
// @param Id query string true "id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/portable/get_new_collar [get]
func GetNewCollarRealtimeByDeviceIdService(deviceId uint) (vo.NewCollar, string){
	accessToken := newCollarAccessToken["accessToken"].(string)
	expireTime := newCollarAccessToken["expireTime"].(time.Time)
	if accessToken == "null" || time.Since(expireTime) >= 24 * time.Hour {
			updateNewCollarToken()
	}
	accessToken = newCollarAccessToken["accessToken"].(string)
	collarDeviceId := dao.GetPortableDeviceInfoById(deviceId).DeviceID
	serverUrl := viper.GetString("newcollar.get-realtime-url")
	payload := url.Values{}
	payload.Set("Iccid", collarDeviceId)
	req, errReq := http.NewRequest("POST", serverUrl, strings.NewReader(payload.Encode()))
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	if errReq != nil {
		panic(errReq.Error())
	}
	req.Header.Set("Authorization", accessToken)
	response, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err.Error())
	}
	code := response.StatusCode
	defer response.Body.Close()
	if code != 200 {
		panic("error occurs during get streaming address from server by deviceID ")
	}
	responseBody := server.GetResponseBodyNewCollarRealtime(response)
	msg := responseBody.Msg
	data := (responseBody.Data)[0]
	newCollarRealtimeData := vo.NewCollar {
		Area: data.Area,
		Iccid: data.Iccid,
		Police: data.Police,
		AllStep: data.AllStep,
		LastTime: data.LastTime,
		Temperature: data.Temperature,
		Station: data.Station,
		IsOnline: data.IsOnline,
		SignalStrength: data.SignalStrength,
		Type: data.Type,
		Voltage: data.Voltage,
	}
	return newCollarRealtimeData, msg
}

// @Summary API of golang gin backend
// @Tags Device-fixed
// @description get all fixed devices by farmhouse : 获取一个牧舍下的所有固定式设备 参数列表：[牧舍ID] 访问携带token
// @version 1.0
// @accept application/json
// @param CompanyId query string true "company id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/fixed/get_by_farmhouse [get]
func GetFixedDeviceListByFarmhouseService(farmhouseId uint) []entity.FixedDevice {
	var fixedDeviceList []entity.FixedDevice
	GetFixedDeviceRecursive(farmhouseId, &fixedDeviceList)
	return fixedDeviceList
}
func GetFixedDeviceRecursive(companyId uint, fixedDeviceList *[]entity.FixedDevice) {
	fixedList := dao.GetFixedDeviceListByFarmhouse(companyId)
	*fixedDeviceList = append(*fixedDeviceList, fixedList...)
	childrenList := dao.GetCompanyListByParent(companyId)
	for _, subCompany := range childrenList {
		GetFixedDeviceRecursive(subCompany.ID, fixedDeviceList)
	}
}

// @Summary API of golang gin backend
// @Tags Device-portable
// @description get all portable devices by farmhouse : 获取一个牧舍下的所有便携式设备 参数列表：[牧舍ID] 访问携带token
// @version 1.0
// @accept application/json
// @param CompanyId query string true "company id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/portable/get_by_farmhouse [get]
func GetPortableDeviceListByFarmhouseService(farmhouseId uint) []vo.BiologyDevice {
	var result []vo.BiologyDevice
	biologyList := GetBiologyListService(farmhouseId)
	for _, biology := range biologyList {
		deviceList := dao.GetPortableDeviceListByBiology(biology.ID)
		for _, device := range deviceList {
			if device.ID != 0 {
				result = append(result, vo.BiologyDevice{
					BiologyId: biology.ID,
					BiologyName: biology.Name,
					BiologyType: biology.BiologyTypeID,
					DeviceId: device.ID,
					DeviceType: device.PortableDeviceTypeID,
				})
			}
		}
	}
	return result
}

// @Summary API of golang gin backend
// @Tags Device-fixed
// @description get user's auth fixed device list : 获取当前用户有权限的所有固定式设备信息 参数列表：[] 访问携带token
// @version 1.0
// @accept application/json
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/fixed/get_auth_list [get]
func GetAuthFixedDeviceListService(userId uint) []vo.AuthFixedDevice {
	companies := dao.GetCompanyListByUserID(userId)
	var childNodeList []uint
	for _, company := range companies {
		getChildNodeRecursive(company.CompanyID, &childNodeList)
	}
	var result []vo.AuthFixedDevice
	for _, node := range childNodeList {
		currList := dao.GetFixedDeviceListByFarmhouse(node)
		for _, curr := range currList {
			result = append(result, vo.AuthFixedDevice{
				DeviceId: curr.ID,
				DeviceType: curr.FixedDeviceTypeID,
				FarmhouseId: curr.FarmhouseID,
				CreateDate: curr.CreatedAt,
				BoughtDate: curr.BoughtTime,
				InstallDate: curr.InstallTime,
				Stat: curr.Stat,
			})
		}
	}
	return result
}

// @Summary API of golang gin backend
// @Tags Device-fixed
// @description get own fixed device list : 获取当前用户拥有的所有固定式设备信息 参数列表：[] 访问携带token
// @version 1.0
// @accept application/json
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/fixed/own_list [get]
func GetOwnFixedDeviceListService(userId uint) []vo.OwnFixedDevice {
	var fixedDeviceInfoList []vo.OwnFixedDevice
	fixedDeviceList := dao.GetOwnFixedDeviceList(userId)
	for _, fixedDevice := range fixedDeviceList {
		fixedDeviceInfoList = append(fixedDeviceInfoList, vo.OwnFixedDevice{
			Id: fixedDevice.ID,
			Type: fixedDevice.FixedDeviceTypeID,
			DeviceId: fixedDevice.DeviceID,
			InstallTime: fixedDevice.InstallTime,
			CreateTime: fixedDevice.CreatedAt,
			FarmhouseId: fixedDevice.FarmhouseID,
			BoughtTime: fixedDevice.BoughtTime,
		})
	}
	return fixedDeviceInfoList
}

// @Summary API of golang gin backend
// @Tags Device-portable
// @description get own portable device list : 获取当前用户拥有的所有便携式设备信息 参数列表：[] 访问携带token
// @version 1.0
// @accept application/json
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/portable/own_list [get]
func GetOwnPortableDeviceListService(userId uint) []vo.OwnPortableDevice {
	var portableDeviceInfoList []vo.OwnPortableDevice
	portableDeviceList := dao.GetOwnPortableDeviceList(userId)
	for _, portableDevice := range portableDeviceList {
		portableDeviceInfoList = append(portableDeviceInfoList, vo.OwnPortableDevice{
			Id: portableDevice.ID,
			Type: portableDevice.PortableDeviceTypeID,
			DeviceId: portableDevice.DeviceID,
			InstallTime: portableDevice.InstallTime,
			CreateTime: portableDevice.CreatedAt,
			BiologyId: portableDevice.BiologyID,
			BoughtTime: portableDevice.BoughtTime,
		})
	}
	return portableDeviceInfoList
}

// @Summary API of golang gin backend
// @Tags Device-portable
// @description get biology portable device list : 获取指定生物绑定的所有便携式设备信息 参数列表：[生物id] 访问携带token
// @version 1.0
// @accept application/json
// @param BiologyId query string true "biology id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/portable/get_by_biology [get]
func GetPortableDeviceListByBiologyService(biologyId uint) []vo.BiologyPortableDevice {
	portableDeviceList := dao.GetPortableDeviceListByBiology(biologyId)
	var result []vo.BiologyPortableDevice
	for _, portableDevice := range portableDeviceList {
		result = append(result, vo.BiologyPortableDevice{
			Id: portableDevice.ID,
			Type: portableDevice.PortableDeviceTypeID,
			DeviceId: portableDevice.DeviceID,
			InstallTime: portableDevice.InstallTime,
			CreateTime: portableDevice.CreatedAt,
			BoughtTime: portableDevice.BoughtTime,
			Stat: portableDevice.Stat,
		})
	}
	return result
}

// @Summary API of golang gin backend
// @Tags Device-Fixed
// @description get fixed device statistic : 获取指定牧场中的固定设备统计信息 参数列表：[牧场Id] 访问携带token
// @version 1.0
// @accept application/json
// @param CompanyId query int true "company id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/fixed/get_statistic [get]
func GetFixedDeviceStatisticService(companyId uint) map[string]uint {
	result := make(map[string]uint)
	deviceList := GetFixedDeviceListByFarmhouseService(companyId)
	for _, device := range deviceList {
		result[device.FixedDeviceTypeID] += 1
	}
	return result
}

// @Summary API of golang gin backend
// @Tags Device-Portable
// @description get portable device statistic : 获取指定牧场中的便携设备统计信息 参数列表：[牧场Id] 访问携带token
// @version 1.0
// @accept application/json
// @param CompanyId query int true "company id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "成功"
// @router /device/portable/get_statistic [get]
func GetPortableDeviceStatisticService(companyId uint) map[string]uint {
	result := make(map[string]uint)
	deviceList := GetPortableDeviceListByFarmhouseService(companyId)
	for _, devcie := range deviceList {
		result[devcie.DeviceType] += 1
	}
	return result
}