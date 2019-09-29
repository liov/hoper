(function() {
    var masterTypeStr = masterType;
    if(masterType != "Plan" && masterType != "PlanStudy" && masterType != "Position" && masterType != "PositionStudy") {
        masterTypeStr = "";
    }
    var arr = '{"OrgID":"' + orgID + '","UserID":"' + userID + '","KnowledgeID":"' + key + '","PackageID":"' + packageID + '","MasterID":"' + masterID + '","MasterType":"' + masterTypeStr + '","PageSize":' + pageSize + ',"StudyTime":0,"studyChapterIDs":"' + stuChpIDs + '","Type":' + type + ',"IsOffLine":false,"DeviceId":"","IsEnd":true,"ReqType":0,"IsCare":true,"viewSchedule":3200}';

    var apitoken = getCookie("ELEARNING_00024");
  var fullUrl = lecaiAPiUrl + "study";
  var dinfineData = {
      url:fullUrl,
      contentType: "application/json",
      dataType: "text",
      type: "POST",
      headers: {
        token: apitoken
      },
      data: arr,
     };
  jQuery.ajax(dinfineData)
})()