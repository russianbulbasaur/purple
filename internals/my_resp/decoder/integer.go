package decoder

import (
	"errors"
	"log"
	"purple/internals/my_resp/constants"
	types "purple/internals/my_resp/purple_data_types"
	"strconv"
)

func decodeInteger(input []byte) (types.PurpleInteger, error, uint64) {
	purpleInteger := types.PurpleInteger{}
	var counter uint64 = 1 //pehla byte
	if !isValid(input) {
		log.Println("Invalid cr lf bytes")
		return purpleInteger, errors.New("invalid cr and lf bytes"), counter
	}
	sign := input[1]
	counter++ //sign byte
	if sign != constants.PositiveSign && sign != constants.NegativeSign {
		log.Println("No sign")
		return purpleInteger, errors.New("no sign"), counter
	}
	integerString := string(input[2:(len(input) - 2)])
	parsedInteger, err := strconv.Atoi(integerString)
	counter += uint64(len(integerString)) + 2
	if err != nil {
		log.Println("Cannot parse the integer value")
		return purpleInteger, errors.New("cannot parse the integer value"), counter
	}
	purpleInteger.Value = parsedInteger
	return purpleInteger, nil, counter
}
