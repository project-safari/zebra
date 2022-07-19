package pkg

// creating some random ports.
func Ports() uint32 {
	nums := []uint32{1, 2, 3, 4, 6, 8, 9, 16, 27, 32, 36, 54, 72, 64, 81, 128, 162, 216, 256, 512}
	length := len(nums)

	ind := RandNum(length)
	port := nums[ind]

	return port
}

// create random sample model names.
func Models() string {
	models := []string{
		"model A", "modelB", "modelC", "modelD",
		"modelA.1", "modelA.1.1", "modelE", "modelD.010",
		"modelF", "modelG", "modelG.1", "modelG.1.1",
	}

	model := RandData(models)

	return model
}

// create random serial codes.
func Serials() string {
	nums := []string{
		"00000", "00001", "00002", "00003",
		"00004", "00005", "00006", "00007",
		"00008", "00009", "00010", "00020",
		"00030", "00040", "00050", "00060",
		"00070", "00080", "00090", "00100",
		"00200", "00300", "00400", "00500",
		"01000", "02000", "03000", "04000",
	}

	ser := RandData(nums)

	return ser
}

// some random rows.
func Rows() string {
	theRows := []string{
		"00A", "00B", "00C", "00D", "00E",
		"0A0", "0B0", "0C0", "0D0", "0E0",
		"A00", "B00", "C00", "D00", "E00",
	}

	oneRow := RandData(theRows)

	return oneRow
}

// create random server IDs.
func SelectServer() string {
	nums := []string{
		"@S00000A", "@S00001B", "@S00002C", "@S00003D",
		"@S00004A", "@S00005B", "@S00006C", "@S00007D",
		"@S00008A", "@S00009B", "@S00010C", "@S00020D",
		"@S00030A", "@S00040B", "@S00050C", "@S00060D",
		"@S00070A", "@S00080B", "@S00090C", "@S00100D",
		"@S00200A", "@S00300B", "@S00400C", "@S00500D",
		"@S01000A", "@S02000B", "@0S3000C", "@S04000D",
	}

	server := RandData(nums)

	return server
}

// create random ESX IDs.
func SelectESX() string {
	nums := []string{
		"@E00000Z", "@E00001Y", "@E00002X", "@E00003W",
		"@E00004Z", "@E00005Y", "@E00006X", "@E00007W",
		"@E00008Z", "@E00009Y", "@E00010X", "@E00020W",
		"@E00030Z", "@E00040Y", "@E00050X", "@E00060W",
		"@E00070Z", "@E00080Y", "@E00090X", "@E00100W",
		"@E00200Z", "@E00300Y", "@E00400X", "@E00500W",
		"@E01000Z", "@E02000Y", "@ES3000X", "@E04000W",
	}

	esx := RandData(nums)

	return esx
}

// create random serial.
func SelectVcenter() string {
	nums := []string{
		"V00000A", "V00001B", "V00002A", "V00003B",
		"V01000B", "V02000A", "V03000B", "V04000A",
		"V00004A", "V00005B", "V00006A", "V00007B",
		"V00008A", "V00009B", "V00010A", "V00020B",
		"V00000B", "V00001A", "V00002B", "V00003A",
		"V00030A", "V00040B", "V00050A", "V00060B",
		"V00004B", "V00005A", "V00006B", "V00007A",
		"V00070A", "V00080B", "V00090A", "V00100B",
		"V00008B", "V00009A", "V00010B", "V00020A",
		"V00200A", "V00300B", "V00400A", "V00500B",
		"V01000A", "V02000B", "V03000A", "V04000B",
		"V00030B", "V00040A", "V00050B", "V00060A",
		"V00200B", "V00300A", "V00400B", "V00500A",
	}

	ser := RandData(nums)

	return ser
}
