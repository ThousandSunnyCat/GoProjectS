package WechatSDK

import (
	"io"
	"encoding/xml"
)

type StringMap map[string]interface{}

type xmlMapEntry struct {
    XMLName xml.Name
    Value   interface{} `xml:",chardata"`
}

func (m StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if len(m) == 0 {
        return nil
    }

    err := e.EncodeToken(start)
    if err != nil {
        return err
    }

    for k, v := range m {
        e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
    }

    return e.EncodeToken(start.End())
}

func (m *StringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    *m = StringMap{}
    for {
        var e xmlMapEntry

        err := d.Decode(&e)
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }

        (*m)[e.XMLName.Local] = e.Value
    }
    return nil
}