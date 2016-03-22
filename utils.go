package main

func CheckErr(err error) {
	if err != nil {
		errl.Println(err)
	}
}
