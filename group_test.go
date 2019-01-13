package main

import "testing"

func TestParseGroup(t *testing.T) {
	group, err := parseGroup("GROUP 5 \"hello there\" 99")
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "group ID", 5, group.ID)
	assertEqualsStr(t, "name", "hello there", group.Name)
	assertEqualsInt(t, "nbr of devices", 1, len(group.Devices))
	assertEqualsInt(t, "device 1 id", 99, group.Devices[0])

	group, err = parseGroup("GROUP     872    \"a\" 7 6 5    88 8  ")
	assertExpectNoErr(t, "", err)
	assertEqualsInt(t, "group ID", 872, group.ID)
	assertEqualsStr(t, "name", "a", group.Name)
	assertEqualsInt(t, "nbr of devices", 5, len(group.Devices))
	assertEqualsInt(t, "device 1 id", 7, group.Devices[0])
	assertEqualsInt(t, "device 2 id", 6, group.Devices[1])
	assertEqualsInt(t, "device 3 id", 5, group.Devices[2])
	assertEqualsInt(t, "device 4 id", 88, group.Devices[3])
	assertEqualsInt(t, "device 5 id", 8, group.Devices[4])
}

func TestGroupToString(t *testing.T) {
	config := "GROUP 5 \"hello there\" 99"
	group, err := parseGroup(config)
	assertExpectNoErr(t, "", err)
	result := group.toConfigStr()
	assertEqualsStr(t, "config", config, result)

	config = "GROUP 1 \"a\""
	group, err = parseGroup(config)
	assertExpectNoErr(t, "", err)
	result = group.toConfigStr()
	assertEqualsStr(t, "config", config, result)

	config = "GROUP 5424 \"hello there or what\" 99 22 43 1 6 7 121"
	group, err = parseGroup(config)
	assertExpectNoErr(t, "", err)
	result = group.toConfigStr()
	assertEqualsStr(t, "config", config, result)

	// Same as above but with whitespace
	group, err = parseGroup("GROUP    5424     \"hello there or what\" 99   22  43 1    6 7  121")
	assertExpectNoErr(t, "", err)
	result = group.toConfigStr()
	assertEqualsStr(t, "config", config, result)
}

func TestParseGroupError(t *testing.T) {
	// Missing quotes
	_, err := parseGroup("GROUP 5 hello 99")
	assertExpectErr(t, "", err)

	// Missing one quote
	_, err = parseGroup("GROUP 5 \"hello 99")
	assertExpectErr(t, "", err)

	// Missing group ID
	_, err = parseGroup("GROUP \"hello there\" 99")
	assertExpectErr(t, "", err)

	// Non numeric group ID
	_, err = parseGroup("GROUP a \"hello there\" 99")
	assertExpectErr(t, "", err)

	// Wrong syntax
	_, err = parseGroup("GROU 5 \"hello\" 99")
	assertExpectErr(t, "", err)

	// Non numeric device ID
	_, err = parseGroup("GROUP 2 \"hello there\" a")
	assertExpectErr(t, "", err)
}
