package superChecker


func ChineseOnly(data interface{})(bool,string,error){
    v :=ToString(data)
    if data == ""{
    	return true,"success",nil
	}
	return true,v,nil
}