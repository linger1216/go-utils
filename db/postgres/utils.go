package postgres

import (
	"bytes"
	"fmt"
	"github.com/linger1216/go-utils/convert"
)

func ValuesPlaceHolder(start int, count int) string {
	var ret bytes.Buffer
	ret.WriteByte('(')
	for i := 1; i <= count; i++ {
		ret.WriteByte('$')
		ret.WriteString(convert.Int64ToString(int64(start + i)))
		if i != count {
			ret.WriteByte(',')
		}
	}
	ret.WriteByte(')')
	return ret.String()
}

func ValuePlaceHolderAndGeometry(start int, count int, geometryPos ...int) string {
	var ret bytes.Buffer
	ret.WriteByte('(')
	for i := 1; i <= count; i++ {
		isGeometry := false
		for _, pos := range geometryPos {
			if i == pos {
				isGeometry = true
				break
			}
		}
		if isGeometry {
			ret.WriteString(fmt.Sprintf(`ST_SetSRID(st_geomfromgeojson($%s),4326)`, convert.Int64ToString(int64(start+i))))
		} else {
			ret.WriteByte('$')
			ret.WriteString(convert.Int64ToString(int64(start + i)))
		}
		if i != count {
			ret.WriteByte(',')
		}
	}
	ret.WriteByte(')')
	return ret.String()
}

func CondSql(first bool) string {
	if first {
		return " where"
	}
	return " and"
}

func SqlStringIn(ids ...string) string {
	var buffer bytes.Buffer
	for _, v := range ids {
		buffer.WriteString("'")
		buffer.WriteString(v)
		buffer.WriteString("'")
		buffer.WriteString(",")
	}
	temp := buffer.String()
	size := len(temp)
	if size > 0 {
		return temp[:size-1]
	}
	return `''`
}

func SqlIntegerIn(ids ...int) string {
	var buffer bytes.Buffer
	for _, v := range ids {
		buffer.WriteString("'")
		buffer.WriteString(convert.Int64ToString(int64(v)))
		buffer.WriteString("'")
		buffer.WriteString(",")
	}
	temp := buffer.String()
	size := len(temp)
	if size > 0 {
		return temp[:size-1]
	}
	return `''`
}

func SqlIntegerArray(ids ...int) string {
	var buffer bytes.Buffer
	buffer.WriteString("array[")
	for i, v := range ids {
		//buffer.WriteString("'")
		buffer.WriteString(convert.Int64ToString(int64(v)))
		//buffer.WriteString("'")
		if i < len(ids)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]::integer[]")
	return buffer.String()
}

func SqlStringArray(ids ...string) string {
	var buffer bytes.Buffer
	buffer.WriteString("array[")
	for i, v := range ids {
		buffer.WriteString("'")
		buffer.WriteString(v)
		buffer.WriteString("'")
		if i < len(ids)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("]::varchar[]")
	return buffer.String()
}

func SqlWithIn(lng, lat float64, radius int) string {
	return fmt.Sprintf(`ST_DWithin(ST_GeomFromText('POINT(%f %f)', 4326)::geography, geometry::geography, %d)`, lng, lat, radius)
}

func ReadTables(uri string) ([]string, error) {
	const ShowTablesQuery = `select relname as table_name,cast(obj_description(relfilenode,'pg_class') as varchar) as comment from pg_class c 
where  relkind = 'r' and relname not like 'pg_%' and relname not like 'sql_%' and relchecks=0 order by relname;`

	db := NewPostgres(NewConfig(uri))

	defer db.Close()

	rows, err := db.Queryx(ShowTablesQuery)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	arr := make([]string, 0)
	for rows.Next() {
		line := make(map[string]interface{})
		err = rows.MapScan(line)
		if err != nil {
			return nil, err
		}

		if v, ok := line["table_name"]; ok {
			t := convert.ToString(v)
			if len(t) > 0 {
				arr = append(arr, t)
			}
		}
	}
	return arr, nil
}
