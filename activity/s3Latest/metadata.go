package s3newestmodel

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
}

type Input struct {
	Bucket     string `md:"bucket,required"`
	Item       string `md:"item,required"`
	File2Check string `md:"file2Check,required"`
	Region     string `md:"region,required"`
	CheckLocal string `md:"checkLocal,required,allowed['file','dir']"`
	CheckS3    string `md:"checkS3,required,allowed['item','prefix']"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	r.Bucket, _ = coerce.ToString(values["bucket"])
	r.Item, _ = coerce.ToString(values["item"])
	r.File2Check, _ = coerce.ToString(values["file2Check"])
	r.Region, _ = coerce.ToString(values["region"])
	r.CheckLocal, _ = coerce.ToString(values["checkLocal"])
	r.CheckS3, _ = coerce.ToString(values["checkS3"])
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"bucket":     r.Bucket,
		"item":       r.Item,
		"file2Check": r.File2Check,
		"region":     r.Region,
		"checkLocal": r.CheckLocal,
		"checkS3":    r.CheckS3,
	}
}

type Output struct {
	ModelFile string `md:"modelFile"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["modelFile"])
	o.ModelFile = strVal
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"modelFile": o.ModelFile,
	}
}
