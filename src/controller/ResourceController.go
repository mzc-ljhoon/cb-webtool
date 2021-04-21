package controller

import (
	// "encoding/json"
	"fmt"
	"github.com/cloud-barista/cb-webtool/src/model"
	service "github.com/cloud-barista/cb-webtool/src/service"
	// util "github.com/cloud-barista/cb-webtool/src/util"

	"github.com/labstack/echo"
	// "io/ioutil"
	"log"
	"net/http"
	//"github.com/davecgh/go-spew/spew"
	echotemplate "github.com/foolin/echo-template"
	echosession "github.com/go-session/echo-session"
)

func ResourceBoard(c echo.Context) error {
	fmt.Println("=========== ResourceBoard start ==============")
	comURL := service.GetCommonURL()
	if loginInfo := service.CallLoginInfo(c); loginInfo.UserID != "" {
		nameSpace := service.GetNameSpaceToString(c)
		fmt.Println("Namespace : ", nameSpace)
		return c.Render(http.StatusOK, "ResourceBoard.html", map[string]interface{}{
			"LoginInfo": loginInfo,
			"NameSpace": nameSpace,
			"comURL":    comURL,
		})
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

// func VpcListForm(c echo.Context) error {
func VpcMngForm(c echo.Context) error {
	fmt.Println("ConnectionConfigList ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	store.Set("namespace", nsList)
	log.Println(" nsList  ", nsList)

	vNetInfoList, respStatus := service.GetVnetList(defaultNameSpaceID)
	// if vNetErr != nil {
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	log.Println("VNetList", vNetInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/NetworkMng", // 파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"CloudOSList":   cloudOsList,
			"NameSpaceList": nsList,
			"VNetList":      vNetInfoList,
			"status":        respStatus.StatusCode,
		})
}

func GetVpcList(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		// return c.JSON(http.StatusBadRequest, map[string]interface{}{
		// 	"message": "invalid tumblebug connection",
		// 	"status":  "403",
		// })
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	vNetInfoList, respStatus := service.GetVnetList(defaultNameSpaceID)
	// if vNetErr != nil {
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":            "success",
		"status":             respStatus.StatusCode,
		"DefaultNameSpaceID": defaultNameSpaceID,
		"VNetList":           vNetInfoList,
	})
}

// Vpc 상세정보
func GetVpcData(c echo.Context) error {
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVNetID := c.Param("vNetID")
	vNetInfo, respStatus := service.GetVpcData(defaultNameSpaceID, paramVNetID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "success",
		"status":   respStatus,
		"VNetInfo": vNetInfo,
	})
}

// Vpc 등록 :
func VpcRegProc(c echo.Context) error {
	log.Println("VpcRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	vNetRegInfo := new(model.VNetRegInfo)
	if err := c.Bind(vNetRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}
	// log.Println(vNetRegInfo)
	resultVNetInfo, respStatus := service.RegVpc(defaultNameSpaceID, vNetRegInfo)
	// respBody, respStatus := service.RegVpc(defaultNameSpaceID, vNetRegInfo)
	// fmt.Println("=============respStatus =============", respStatus)
	// fmt.Println("=============respBody ===============", respBody)

	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "success",
		"status":   respStatus.StatusCode,
		"VNetInfo": resultVNetInfo,
	})
}

// 삭제
func VpcDelProc(c echo.Context) error {
	log.Println("ConnectionConfigDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVNetID := c.Param("vNetID")

	respBody, respStatus := service.DelVpc(defaultNameSpaceID, paramVNetID)
	fmt.Println("=============respBody =============", respBody)

	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

// SecurityGroup 관리 화면
func SecirityGroupMngForm(c echo.Context) error {
	fmt.Println("SecirityGroupMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	store.Set("namespace", nsList)
	log.Println(" nsList  ", nsList)

	securityGroupInfoList, respStatus := service.GetSecurityGroupList(defaultNameSpaceID)
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	log.Println("securityGroupInfoList", securityGroupInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/SecurityGroupMng", // 파일명
		map[string]interface{}{
			"LoginInfo":         loginInfo,
			"CloudOSList":       cloudOsList,
			"NameSpaceList":     nsList,
			"SecurityGroupList": securityGroupInfoList,
			"status":            respStatus.StatusCode,
		})
}

// SecurityGroup 목록
func GetSecirityGroupList(c echo.Context) error {
	log.Println("GetSecirityGroupList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	securityGroupInfoList, respStatus := service.GetSecurityGroupList(defaultNameSpaceID)
	// if vNetErr != nil {
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":            "success",
		"status":             respStatus.StatusCode,
		"DefaultNameSpaceID": defaultNameSpaceID,
		"SecurityGroupList":  securityGroupInfoList,
	})
}

// 상세정보
func GetSecirityGroupData(c echo.Context) error {
	log.Println("GetSecirityGroupData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSecurityGroupID := c.Param("securityGroupID")
	securityGroupInfo, respStatus := service.GetSecurityGroupData(defaultNameSpaceID, paramSecurityGroupID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":           "success",
		"status":            respStatus,
		"SecurityGroupInfo": securityGroupInfo,
	})
}

// 등록 :
func SecirityGroupRegProc(c echo.Context) error {
	log.Println("SecirityGroupRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	securityGroupRegInfo := new(model.SecurityGroupRegInfo)
	if err := c.Bind(securityGroupRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	resultSecurityGroupInfo, respStatus := service.RegSecurityGroup(defaultNameSpaceID, securityGroupRegInfo)

	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":           "success",
		"status":            respStatus.StatusCode,
		"SecurityGroupInfo": resultSecurityGroupInfo,
	})
}

// 삭제
func SecirityGroupDelProc(c echo.Context) error {
	log.Println("SecirityGroupDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSecurityGroupID := c.Param("securityGroupID")

	respBody, respStatus := service.DelSecurityGroup(defaultNameSpaceID, paramSecurityGroupID)
	fmt.Println("=============respBody =============", respBody)

	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

func SshKeyMngForm(c echo.Context) error {
	fmt.Println("SshKeyMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	store.Set("namespace", nsList)
	log.Println(" nsList  ", nsList)

	sshKeyInfoList, respStatus := service.GetSshKeyInfoList(defaultNameSpaceID)
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	log.Println("sshKeyInfoList", sshKeyInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/SshKeyMng", // 파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"CloudOSList":   cloudOsList,
			"NameSpaceList": nsList,
			"SshKeyList":    sshKeyInfoList,
			"status":        respStatus.StatusCode,
		})
}

func GetSshKeyList(c echo.Context) error {
	log.Println("GetSshKeyList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	sshKeyInfoList, respStatus := service.GetSshKeyInfoList(defaultNameSpaceID)
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":            "success",
		"status":             respStatus.StatusCode,
		"DefaultNameSpaceID": defaultNameSpaceID,
		"SshKeyList":         sshKeyInfoList,
	})
}

// SSHKey 상세정보
func GetSshKeyData(c echo.Context) error {
	log.Println("GetSshKeyData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSshKey := c.Param("sshKeyID")
	sshKeyInfo, respStatus := service.GetSshKeyData(defaultNameSpaceID, paramSshKey)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus,
		"SshKeyInfo": sshKeyInfo,
	})
}

// SSHKey 등록 :
func SshKeyRegProc(c echo.Context) error {
	log.Println("SshKeyRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	sshKeyRegInfo := new(model.SshKeyRegInfo)
	if err := c.Bind(sshKeyRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	resultSshKeyInfo, respStatus := service.RegSshKey(defaultNameSpaceID, sshKeyRegInfo)

	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"status":     respStatus.StatusCode,
		"SshKeyInfo": resultSshKeyInfo,
	})
}

// 삭제
func SshKeyDelProc(c echo.Context) error {
	log.Println("SshKeyDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramSshKeyID := c.Param("sshKeyID")

	respBody, respStatus := service.DelSshKey(defaultNameSpaceID, paramSshKeyID)
	fmt.Println("=============respBody =============", respBody)

	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

// VirtualMachine Image 등록 form
func VirtualMachineImageMngForm(c echo.Context) error {
	fmt.Println("VirtualMachineImageMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	store.Set("namespace", nsList)
	log.Println(" nsList  ", nsList)

	virtualMachineImageInfoList, respStatus := service.GetVirtualMachineImageInfoList(defaultNameSpaceID)
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	log.Println("VirtualMachineImageInfoList", virtualMachineImageInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/VirtualMachineImageMng", // 파일명
		map[string]interface{}{
			"LoginInfo":               loginInfo,
			"CloudOSList":             cloudOsList,
			"NameSpaceList":           nsList,
			"VirtualMachineImageList": virtualMachineImageInfoList,
			"status":                  respStatus.StatusCode,
		})
}

func GetVirtualMachineImageList(c echo.Context) error {
	log.Println("GetVirtualMachineImageList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	virtualMachineImageInfoList, respStatus := service.GetVirtualMachineImageInfoList(defaultNameSpaceID)
	// if vNetErr != nil {
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                 "success",
		"status":                  respStatus.StatusCode,
		"DefaultNameSpaceID":      defaultNameSpaceID,
		"VirtualMachineImageList": virtualMachineImageInfoList,
	})
}

// VirtualMachineImage 상세정보
func GetVirtualMachineImageData(c echo.Context) error {
	log.Println("GetVirtualMachineImageData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVirtualMachineImage := c.Param("imageID")
	virtualMachineImageInfo, respStatus := service.GetVirtualMachineImageData(defaultNameSpaceID, paramVirtualMachineImage)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                 "success",
		"status":                  respStatus,
		"VirtualMachineImageInfo": virtualMachineImageInfo,
	})
}

// VirtualMachineImage 등록 :
func VirtualMachineImageRegProc(c echo.Context) error {
	log.Println("VirtualMachineImageRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	virtualMachineImageRegInfo := new(model.VirtualMachineImageRegInfo)
	if err := c.Bind(virtualMachineImageRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	resultVirtualMachineImageInfo, respStatus := service.RegVirtualMachineImage(defaultNameSpaceID, virtualMachineImageRegInfo)
	// todo : return message 조치 필요. 중복 등 에러났을 때 message 표시가 제대로 되지 않음
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	// respBody := resp.Body
	// respStatusCode := resp.StatusCode
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                 "success",
		"status":                  respStatus.StatusCode,
		"VirtualMachineImageInfo": resultVirtualMachineImageInfo,
	})
}

// 삭제
func VirtualMachineImageDelProc(c echo.Context) error {
	log.Println("VirtualMachineImageDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVirtualMachineImageID := c.Param("imageID")

	respBody, respStatus := service.DelVirtualMachineImage(defaultNameSpaceID, paramVirtualMachineImageID)
	fmt.Println("=============respBody =============", respBody)

	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

// 해당 namespace의 모든 VMImage 삭제
func AllVirtualMachineImageDelProc(c echo.Context) error {
	log.Println("AllVirtualMachineImageDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramNameSpaceID := c.Param("nameSpaceID")

	// 해당 Namespace의 모든 Image가 삭제 되므로 선택한 namespace와 defaultNamespace가 같아야 삭제가능
	if defaultNameSpaceID != paramNameSpaceID {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "현재 namespace만 삭제 가능합니다.",
			"status":  4040,
		})
	}

	respBody, respStatus := service.DelAllVirtualMachineImage(defaultNameSpaceID)
	fmt.Println("=============respBody =============", respBody)

	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

// connection에 해당하는 machine image 목록
// connection이 resion별로 생성되므로 결국 해당 provider의 resion 내 vm목록을 가져옴
func LookupVirtualMachineImageList(c echo.Context) error {
	log.Println("GetVirtualMachineImageList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// paramConnectionName := c.Param("connectionName")
	paramConnectionName := c.QueryParam("connectionName")

	log.Println("paramConnectionName : ", paramConnectionName)
	virtualMachineImageInfoList, respStatus := service.LookupVirtualMachineImageList(paramConnectionName)
	// if vNetErr != nil {
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
		// "DefaultNameSpaceID": defaultNameSpaceID,
		"VirtualMachineImageList": virtualMachineImageInfoList,
	})
}

// lookupImage 상세정보
func LookupVirtualMachineImageData(c echo.Context) error {
	log.Println("LookupVirtualMachineImageData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramVirtualMachineImage := c.Param("imageID")
	virtualMachineImageInfo, respStatus := service.LookupVirtualMachineImageData(paramVirtualMachineImage)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                 "success",
		"status":                  respStatus,
		"VirtualMachineImageInfo": virtualMachineImageInfo,
	})
}

// lookupImage 상세정보
func SearchVirtualMachineImageList(c echo.Context) error {
	log.Println("SearchVirtualMachineImageList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramKeywords := c.Param("keywords")
	virtualMachineImageInfoList, respStatus := service.SearchVirtualMachineImageList(defaultNameSpaceID, paramKeywords)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                 "success",
		"status":                  respStatus,
		"VirtualMachineImageList": virtualMachineImageInfoList,
	})
}

// TODO : Fetch 의 의미 파악
func FetchVirtualMachineImageList(c echo.Context) error {
	log.Println("FetchVirtualMachineImageList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	virtualMachineImageInfoList, respStatus := service.FetchVirtualMachineImageList(defaultNameSpaceID)
	// if vNetErr != nil {
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
		// "DefaultNameSpaceID": defaultNameSpaceID,
		"VirtualMachineImageList": virtualMachineImageInfoList,
	})
}

// VMSpecMng 등록 form
func VmSpecMngForm(c echo.Context) error {
	fmt.Println("VmSpecMngForm ************ : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	store := echosession.FromContext(c)

	cloudOsList, _ := service.GetCloudOSList()
	store.Set("cloudos", cloudOsList)
	log.Println(" cloudOsList  ", cloudOsList)

	// 최신 namespacelist 가져오기
	nsList, _ := service.GetNameSpaceList()
	store.Set("namespace", nsList)
	log.Println(" nsList  ", nsList)

	vmSpecInfoList, respStatus := service.GetVmSpecInfoList(defaultNameSpaceID)
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	log.Println("VmSpecInfoList", vmSpecInfoList)

	return echotemplate.Render(c, http.StatusOK,
		"setting/resources/VirtualMachineSpecMng", // 파일명
		map[string]interface{}{
			"LoginInfo":     loginInfo,
			"CloudOSList":   cloudOsList,
			"NameSpaceList": nsList,
			"VmSpecList":    vmSpecInfoList,
			"status":        respStatus.StatusCode,
		})
}

func GetVmSpecList(c echo.Context) error {
	log.Println("GetVmSpecList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	vmSpecInfoList, respStatus := service.GetVmSpecInfoList(defaultNameSpaceID)
	// if vNetErr != nil {
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":            "success",
		"status":             respStatus.StatusCode,
		"DefaultNameSpaceID": defaultNameSpaceID,
		"VmSpecList":         vmSpecInfoList,
	})
}

// VMSpec 상세정보
func GetVmSpecData(c echo.Context) error {
	log.Println("GetVmSpecData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVMSpecID := c.Param("specID")
	vmSpecInfo, respStatus := service.GetVmSpecInfoData(defaultNameSpaceID, paramVMSpecID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"VmSpec":  vmSpecInfo,
	})
}

// VMSpec 등록 :
func VmSpecRegProc(c echo.Context) error {
	log.Println("VMSpecRegProc : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	vmSpecRegInfo := new(model.VmSpecRegInfo)
	if err := c.Bind(vmSpecRegInfo); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "fail",
			"status":  "fail",
		})
	}

	resultVmSpecInfo, respStatus := service.RegVmSpec(defaultNameSpaceID, vmSpecRegInfo)
	// todo : return message 조치 필요. 중복 등 에러났을 때 message 표시가 제대로 되지 않음
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}
	// respBody := resp.Body
	// respStatusCode := resp.StatusCode
	// respStatus := resp.Status
	// log.Println("respStatusCode = ", respStatusCode)
	// log.Println("respStatus = ", respStatus)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
		"VMSpec":  resultVmSpecInfo,
	})
}

// 삭제
func VmSpecDelProc(c echo.Context) error {
	log.Println("VmSpecDelProc : ")

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	paramVMSpecID := c.Param("specID")

	respBody, respStatus := service.DelVMSpec(defaultNameSpaceID, paramVMSpecID)
	fmt.Println("=============respBody =============", respBody)

	// if reErr != nil {
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
	})
}

// lookupImage 목록
func LookupVmSpecList(c echo.Context) error {
	log.Println("LookupVmSpecList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// store := echosession.FromContext(c)
	defaultNameSpaceID := loginInfo.DefaultNameSpaceID
	// TODO : defaultNameSpaceID 가 없으면 설정화면으로 보낼 것
	vmSpecInfoList, respStatus := service.LookupVmSpecInfoList()
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":            "success",
		"status":             respStatus.StatusCode,
		"DefaultNameSpaceID": defaultNameSpaceID,
		"VmSpecList":         vmSpecInfoList,
	})
}

// lookupImage 상세정보
func LookupVmSpecData(c echo.Context) error {
	log.Println("LookupVmSpecData : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	paramVMSpec := c.Param("specID") // IID
	vmSpecInfo, respStatus := service.LookupVmSpecInfoData(paramVMSpec)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus,
		"VmSpec":  vmSpecInfo,
	})
}

// TODO : Fetch 의미 파악필요
func FetchVmSpecList(c echo.Context) error {
	log.Println("FetchVMSpecList : ")
	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	defaultNameSpaceID := loginInfo.DefaultNameSpaceID

	vmSpecInfoList, respStatus := service.FetchVmSpecInfoList(defaultNameSpaceID)
	if respStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": respStatus.Message,
			"status":  respStatus.StatusCode,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"status":  respStatus.StatusCode,
		// "DefaultNameSpaceID": defaultNameSpaceID,
		"VmSpec": vmSpecInfoList,
	})
}

// resourcesGroup.PUT("/vmspec/put/:specID", controller.VMSpecPutProc)	// RegProc _ SshKey 같이 앞으로 넘길까
// resourcesGroup.POST("/vmspec/filterspecs", controller.FilterVMSpecList)
// resourcesGroup.POST("/vmspec/filterspecsbyrange", controller.FilterVMSpecListByRange)
