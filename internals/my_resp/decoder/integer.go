package decoder

//func decodeInteger(input []byte) (constants.PurpleInteger, error) {
//	if !isValid(input) {
//		log.Println("Invalid cr lf bytes")
//		return -1, errors.New("invalid cr and lf bytes")
//	}
//	sign := input[1]
//	if sign != constants.PositiveSign && sign != constants.NegativeSign {
//		log.Println("No sign")
//		return -1, errors.New("no sign")
//	}
//	parsedInteger, err := strconv.Atoi(string(input[2:(len(input) - 2)]))
//	if err != nil {
//		log.Println("Cannot parse the integer value")
//		return -1, errors.New("cannot parse the integer value")
//	}
//	return constants.PurpleInteger(parsedInteger), nil
//}
