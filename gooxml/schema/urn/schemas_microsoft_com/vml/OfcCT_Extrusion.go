// Copyright 2017 Baliance. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased by contacting sales@baliance.com.

package vml

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"baliance.com/gooxml/schema/soo/ofc/sharedTypes"
)

type OfcCT_Extrusion struct {
	OnAttr                 sharedTypes.ST_TrueFalse
	TypeAttr               OfcST_ExtrusionType
	RenderAttr             OfcST_ExtrusionRender
	ViewpointoriginAttr    *string
	ViewpointAttr          *string
	PlaneAttr              OfcST_ExtrusionPlane
	SkewangleAttr          *float32
	SkewamtAttr            *string
	ForedepthAttr          *string
	BackdepthAttr          *string
	OrientationAttr        *string
	OrientationangleAttr   *float32
	LockrotationcenterAttr sharedTypes.ST_TrueFalse
	AutorotationcenterAttr sharedTypes.ST_TrueFalse
	RotationcenterAttr     *string
	RotationangleAttr      *string
	ColormodeAttr          OfcST_ColorMode
	ColorAttr              *string
	ShininessAttr          *float32
	SpecularityAttr        *string
	DiffusityAttr          *string
	MetalAttr              sharedTypes.ST_TrueFalse
	EdgeAttr               *string
	FacetAttr              *string
	LightfaceAttr          sharedTypes.ST_TrueFalse
	BrightnessAttr         *string
	LightpositionAttr      *string
	LightlevelAttr         *string
	LightharshAttr         sharedTypes.ST_TrueFalse
	Lightposition2Attr     *string
	Lightlevel2Attr        *string
	Lightharsh2Attr        sharedTypes.ST_TrueFalse
	ExtAttr                ST_Ext
}

func NewOfcCT_Extrusion() *OfcCT_Extrusion {
	ret := &OfcCT_Extrusion{}
	return ret
}

func (m *OfcCT_Extrusion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if m.OnAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.OnAttr.MarshalXMLAttr(xml.Name{Local: "on"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.TypeAttr != OfcST_ExtrusionTypeUnset {
		attr, err := m.TypeAttr.MarshalXMLAttr(xml.Name{Local: "type"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.RenderAttr != OfcST_ExtrusionRenderUnset {
		attr, err := m.RenderAttr.MarshalXMLAttr(xml.Name{Local: "render"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.ViewpointoriginAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "viewpointorigin"},
			Value: fmt.Sprintf("%v", *m.ViewpointoriginAttr)})
	}
	if m.ViewpointAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "viewpoint"},
			Value: fmt.Sprintf("%v", *m.ViewpointAttr)})
	}
	if m.PlaneAttr != OfcST_ExtrusionPlaneUnset {
		attr, err := m.PlaneAttr.MarshalXMLAttr(xml.Name{Local: "plane"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.SkewangleAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "skewangle"},
			Value: fmt.Sprintf("%v", *m.SkewangleAttr)})
	}
	if m.SkewamtAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "skewamt"},
			Value: fmt.Sprintf("%v", *m.SkewamtAttr)})
	}
	if m.ForedepthAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "foredepth"},
			Value: fmt.Sprintf("%v", *m.ForedepthAttr)})
	}
	if m.BackdepthAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "backdepth"},
			Value: fmt.Sprintf("%v", *m.BackdepthAttr)})
	}
	if m.OrientationAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "orientation"},
			Value: fmt.Sprintf("%v", *m.OrientationAttr)})
	}
	if m.OrientationangleAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "orientationangle"},
			Value: fmt.Sprintf("%v", *m.OrientationangleAttr)})
	}
	if m.LockrotationcenterAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.LockrotationcenterAttr.MarshalXMLAttr(xml.Name{Local: "lockrotationcenter"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.AutorotationcenterAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.AutorotationcenterAttr.MarshalXMLAttr(xml.Name{Local: "autorotationcenter"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.RotationcenterAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "rotationcenter"},
			Value: fmt.Sprintf("%v", *m.RotationcenterAttr)})
	}
	if m.RotationangleAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "rotationangle"},
			Value: fmt.Sprintf("%v", *m.RotationangleAttr)})
	}
	if m.ColormodeAttr != OfcST_ColorModeUnset {
		attr, err := m.ColormodeAttr.MarshalXMLAttr(xml.Name{Local: "colormode"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.ColorAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "color"},
			Value: fmt.Sprintf("%v", *m.ColorAttr)})
	}
	if m.ShininessAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "shininess"},
			Value: fmt.Sprintf("%v", *m.ShininessAttr)})
	}
	if m.SpecularityAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "specularity"},
			Value: fmt.Sprintf("%v", *m.SpecularityAttr)})
	}
	if m.DiffusityAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "diffusity"},
			Value: fmt.Sprintf("%v", *m.DiffusityAttr)})
	}
	if m.MetalAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.MetalAttr.MarshalXMLAttr(xml.Name{Local: "metal"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.EdgeAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "edge"},
			Value: fmt.Sprintf("%v", *m.EdgeAttr)})
	}
	if m.FacetAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "facet"},
			Value: fmt.Sprintf("%v", *m.FacetAttr)})
	}
	if m.LightfaceAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.LightfaceAttr.MarshalXMLAttr(xml.Name{Local: "lightface"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.BrightnessAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "brightness"},
			Value: fmt.Sprintf("%v", *m.BrightnessAttr)})
	}
	if m.LightpositionAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "lightposition"},
			Value: fmt.Sprintf("%v", *m.LightpositionAttr)})
	}
	if m.LightlevelAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "lightlevel"},
			Value: fmt.Sprintf("%v", *m.LightlevelAttr)})
	}
	if m.LightharshAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.LightharshAttr.MarshalXMLAttr(xml.Name{Local: "lightharsh"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.Lightposition2Attr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "lightposition2"},
			Value: fmt.Sprintf("%v", *m.Lightposition2Attr)})
	}
	if m.Lightlevel2Attr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "lightlevel2"},
			Value: fmt.Sprintf("%v", *m.Lightlevel2Attr)})
	}
	if m.Lightharsh2Attr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.Lightharsh2Attr.MarshalXMLAttr(xml.Name{Local: "lightharsh2"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.ExtAttr != ST_ExtUnset {
		attr, err := m.ExtAttr.MarshalXMLAttr(xml.Name{Local: "ext"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	e.EncodeToken(start)
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *OfcCT_Extrusion) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Local == "colormode" {
			m.ColormodeAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "color" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.ColorAttr = &parsed
			continue
		}
		if attr.Name.Local == "type" {
			m.TypeAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "shininess" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			pt := float32(parsed)
			m.ShininessAttr = &pt
			continue
		}
		if attr.Name.Local == "viewpointorigin" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.ViewpointoriginAttr = &parsed
			continue
		}
		if attr.Name.Local == "specularity" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.SpecularityAttr = &parsed
			continue
		}
		if attr.Name.Local == "plane" {
			m.PlaneAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "diffusity" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.DiffusityAttr = &parsed
			continue
		}
		if attr.Name.Local == "skewamt" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.SkewamtAttr = &parsed
			continue
		}
		if attr.Name.Local == "metal" {
			m.MetalAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "backdepth" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.BackdepthAttr = &parsed
			continue
		}
		if attr.Name.Local == "edge" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.EdgeAttr = &parsed
			continue
		}
		if attr.Name.Local == "lightlevel2" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.Lightlevel2Attr = &parsed
			continue
		}
		if attr.Name.Local == "orientationangle" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			pt := float32(parsed)
			m.OrientationangleAttr = &pt
			continue
		}
		if attr.Name.Local == "on" {
			m.OnAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "lightharsh" {
			m.LightharshAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "lightface" {
			m.LightfaceAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "foredepth" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.ForedepthAttr = &parsed
			continue
		}
		if attr.Name.Local == "ext" {
			m.ExtAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "autorotationcenter" {
			m.AutorotationcenterAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "facet" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.FacetAttr = &parsed
			continue
		}
		if attr.Name.Local == "render" {
			m.RenderAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "lightlevel" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.LightlevelAttr = &parsed
			continue
		}
		if attr.Name.Local == "brightness" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.BrightnessAttr = &parsed
			continue
		}
		if attr.Name.Local == "skewangle" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			pt := float32(parsed)
			m.SkewangleAttr = &pt
			continue
		}
		if attr.Name.Local == "lightposition2" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.Lightposition2Attr = &parsed
			continue
		}
		if attr.Name.Local == "rotationangle" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.RotationangleAttr = &parsed
			continue
		}
		if attr.Name.Local == "lightharsh2" {
			m.Lightharsh2Attr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "orientation" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.OrientationAttr = &parsed
			continue
		}
		if attr.Name.Local == "lockrotationcenter" {
			m.LockrotationcenterAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "rotationcenter" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.RotationcenterAttr = &parsed
			continue
		}
		if attr.Name.Local == "viewpoint" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.ViewpointAttr = &parsed
			continue
		}
		if attr.Name.Local == "lightposition" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.LightpositionAttr = &parsed
			continue
		}
	}
	// skip any extensions we may find, but don't support
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("parsing OfcCT_Extrusion: %s", err)
		}
		if el, ok := tok.(xml.EndElement); ok && el.Name == start.Name {
			break
		}
	}
	return nil
}

// Validate validates the OfcCT_Extrusion and its children
func (m *OfcCT_Extrusion) Validate() error {
	return m.ValidateWithPath("OfcCT_Extrusion")
}

// ValidateWithPath validates the OfcCT_Extrusion and its children, prefixing error messages with path
func (m *OfcCT_Extrusion) ValidateWithPath(path string) error {
	if err := m.OnAttr.ValidateWithPath(path + "/OnAttr"); err != nil {
		return err
	}
	if err := m.TypeAttr.ValidateWithPath(path + "/TypeAttr"); err != nil {
		return err
	}
	if err := m.RenderAttr.ValidateWithPath(path + "/RenderAttr"); err != nil {
		return err
	}
	if err := m.PlaneAttr.ValidateWithPath(path + "/PlaneAttr"); err != nil {
		return err
	}
	if err := m.LockrotationcenterAttr.ValidateWithPath(path + "/LockrotationcenterAttr"); err != nil {
		return err
	}
	if err := m.AutorotationcenterAttr.ValidateWithPath(path + "/AutorotationcenterAttr"); err != nil {
		return err
	}
	if err := m.ColormodeAttr.ValidateWithPath(path + "/ColormodeAttr"); err != nil {
		return err
	}
	if err := m.MetalAttr.ValidateWithPath(path + "/MetalAttr"); err != nil {
		return err
	}
	if err := m.LightfaceAttr.ValidateWithPath(path + "/LightfaceAttr"); err != nil {
		return err
	}
	if err := m.LightharshAttr.ValidateWithPath(path + "/LightharshAttr"); err != nil {
		return err
	}
	if err := m.Lightharsh2Attr.ValidateWithPath(path + "/Lightharsh2Attr"); err != nil {
		return err
	}
	if err := m.ExtAttr.ValidateWithPath(path + "/ExtAttr"); err != nil {
		return err
	}
	return nil
}
