package utils

import (
	"fmt"
	"strings"
)

func EncodeCommand(parts ...string) string {
	/*
		   RESP(Redis protocol) structure is like this

		   In this case our message is array of parts

		   where length is 3 hence *3

		   and \r\n	returns to the beginning and moves down to the next line.

		   its the delimiter

		   *3\r\n
		   $7\r\nPUBLISH\r\n
		   $6\r\ntokens\r\n
		   $32\r\n<your token>\r\n

		   where $7 represents the length of string
		   in this case P U B L I S H
		   				1 2 3 4 5 6 7
			then next line and similer impelentation
	*/
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("*%d\r\n", len(parts)))

	for _, p := range parts {
		sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(p), p))
	}

	return sb.String()
}
