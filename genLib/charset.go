// charset.go

/// +build OMIT
package genLib

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/net/html/charset"
	iconv "gopkg.in/iconv.v1"
)

// Detect encoding type for strings
func DetectCharsetStr(inputStr string) (name string) {
	_, name, _ = charset.DetermineEncoding([]byte(inputStr), "")
	return strings.ToLower(name)
}

// Detect encoding type for files
func DetectCharsetFile(filename string) (name string) {
	textFileBytes, err := ioutil.ReadFile(filename)
	Check(err, `DetectCharsetFile: `+filename)
	_, name, _ = charset.DetermineEncoding(textFileBytes, "")
	return strings.ToLower(name)
}

// Get charset and convert to utf-8 if needed
func CharsetToUtf8(str string) string {
	charSET := DetectCharsetStr(str)
	if charSET != "utf-8" {
		r, _ := charset.NewReader(strings.NewReader(str), charSET) // convert to UTF-8
		strByte, _ := ioutil.ReadAll(r)
		return string(strByte)
	}
	return str
}

// Change charset of a text file
func FileCharsetSwitch(inCharset, outCharset, inFilename string, outFilename ...string) {
	var flagReplace bool
	if len(outFilename) == 0 {
		outFilename = append(outFilename, inFilename+".~~~")
		flagReplace = true
	} else if inFilename == outFilename[0] {
		outFilename = append(outFilename, inFilename+".~~~")
		flagReplace = true
	}

	cd, err := iconv.Open(outCharset, inCharset)
	Check(err, "iconv.Open failed!")
	defer cd.Close()

	osOpen, err := os.Open(inFilename)
	Check(err, "os.Open failed!")

	r := iconv.NewReader(cd, osOpen, 0)

	osWrite, err := os.Create(outFilename[0])
	defer osWrite.Close()
	Check(err, "os.Create failed!")
	_, err = io.Copy(osWrite, r)
	Check(err, "io.Copy failed!")

	if flagReplace {
		err = os.Remove(inFilename)
		Check(err, "os.Remove failed!")
		err = os.Rename(outFilename[0], inFilename)
		Check(err, "os.Rename failed!")
	}
}

type CharsetList struct {
	SimpleCharsetList []string
	FullCharsetList   []string
}

func NewCharsetList() CharsetList {
	cl := CharsetList{}
	cl.Init()
	return cl
}

func (cl *CharsetList) GetPos(ch string) (posLong, posShort int) {
	ch = strings.ToLower(ch)
	posLong = -1
	posShort = -1
	for idx, value := range cl.FullCharsetList {
		if value == ch {
			posLong = idx
			break
		}
	}

	for idx, value := range cl.SimpleCharsetList {
		if value == ch {
			posShort = idx
			break
		}
	}
	return posLong, posShort
}

func (cl *CharsetList) GetCharset(pos int) (chLong, chShort string) {
	if len(cl.FullCharsetList)-1 > pos {
		chLong = cl.FullCharsetList[pos]
	} else {
		chLong = ""
	}
	if len(cl.SimpleCharsetList)-1 > pos {
		chShort = cl.SimpleCharsetList[pos]
	} else {
		chShort = ""
	}
	return chLong, chShort
}

func (cl *CharsetList) Init() {
	cl.SimpleCharsetList = []string{"ascii", "big-5", "euc-cn", "euc-jisx0213", "euc-jp-ms", "euc-jp", "euc-kr", "euc-tw",
		"euccn", "eucjp-ms", "eucjp-open", "eucjp-win", "eucjp", "euckr", "gb", "gb2312", "gb13000", "gb18030", "gbk", "iso-8859-1",
		"iso-8859-2", "iso-8859-3", "iso-8859-4", "iso-8859-5", "iso-8859-6", "iso-8859-7", "iso-8859-8", "iso-8859-9", "iso-8859-9e",
		"iso-8859-10", "iso-8859-11", "iso-8859-13", "iso-8859-14", "iso-8859-15", "iso-8859-16", "shift-jis", "utf-8", "utf-16",
		"utf-16be", "utf-16le", "utf-32", "utf-32be", "utf-32le", "utf7", "utf8", "utf16", "utf16be", "utf16le", "utf32", "utf32be",
		"utf32le", "windows-31j", "windows-874", "windows-936", "windows-1250", "windows-1251", "windows-1252", "windows-1253",
		"windows-1254", "windows-1255", "windows-1256", "windows-1257", "windows-1258"}
	cl.FullCharsetList = []string{"437", "500", "500v1", "850", "851", "852", "855", "856", "857",
		"860", "861", "862", "863", "864", "865", "866", "866nav", "869", "874", "904", "1026", "1046",
		"1047", "8859_1", "8859_2", "8859_3", "8859_4", "8859_5", "8859_6", "8859_7", "8859_8", "8859_9",
		"10646-1:1993", "10646-1:1993/ucs4", "ansi_x3.4-1968", "ansi_x3.4-1986", "ansi_x3.4", "ansi_x3.110-1983",
		"ansi_x3.110", "arabic", "arabic7", "armscii-8", "ascii", "asmo-708", "asmo_449", "baltic", "big-5",
		"big-five", "big5-hkscs", "big5", "big5hkscs", "bigfive", "brf", "bs_4730", "ca", "cn-big5",
		"cn-gb", "cn", "cp-ar", "cp-gr", "cp-hu", "cp037", "cp038", "cp273", "cp274", "cp275", "cp278",
		"cp280", "cp281", "cp282", "cp284", "cp285", "cp290", "cp297", "cp367", "cp420", "cp423", "cp424",
		"cp437", "cp500", "cp737", "cp770", "cp771", "cp772", "cp773", "cp774", "cp775", "cp803", "cp813",
		"cp819", "cp850", "cp851", "cp852", "cp855", "cp856", "cp857", "cp860", "cp861", "cp862", "cp863",
		"cp864", "cp865", "cp866", "cp866nav", "cp868", "cp869", "cp870", "cp871", "cp874", "cp875",
		"cp880", "cp891", "cp901", "cp902", "cp903", "cp904", "cp905", "cp912", "cp915", "cp916", "cp918",
		"cp920", "cp921", "cp922", "cp930", "cp932", "cp933", "cp935", "cp936", "cp937", "cp939", "cp949",
		"cp950", "cp1004", "cp1008", "cp1025", "cp1026", "cp1046", "cp1047", "cp1070", "cp1079",
		"cp1081", "cp1084", "cp1089", "cp1097", "cp1112", "cp1122", "cp1123", "cp1124", "cp1125",
		"cp1129", "cp1130", "cp1132", "cp1133", "cp1137", "cp1140", "cp1141", "cp1142", "cp1143",
		"cp1144", "cp1145", "cp1146", "cp1147", "cp1148", "cp1149", "cp1153", "cp1154", "cp1155",
		"cp1156", "cp1157", "cp1158", "cp1160", "cp1161", "cp1162", "cp1163", "cp1164", "cp1166",
		"cp1167", "cp1250", "cp1251", "cp1252", "cp1253", "cp1254", "cp1255", "cp1256", "cp1257",
		"cp1258", "cp1282", "cp1361", "cp1364", "cp1371", "cp1388", "cp1390", "cp1399", "cp4517",
		"cp4899", "cp4909", "cp4971", "cp5347", "cp9030", "cp9066", "cp9448", "cp10007", "cp12712",
		"cp16804", "cpibm861", "csa7-1", "csa7-2", "csascii", "csa_t500-1983", "csa_t500",
		"csa_z243.4-1985-1", "csa_z243.4-1985-2", "csa_z243.419851", "csa_z243.419852",
		"csdecmcs", "csebcdicatde", "csebcdicatdea", "csebcdiccafr", "csebcdicdkno",
		"csebcdicdknoa", "csebcdices", "csebcdicesa", "csebcdicess", "csebcdicfise",
		"csebcdicfisea", "csebcdicfr", "csebcdicit", "csebcdicpt", "csebcdicuk", "csebcdicus",
		"cseuckr", "cseucpkdfmtjapanese", "csgb2312", "cshproman8", "csibm037", "csibm038",
		"csibm273", "csibm274", "csibm275", "csibm277", "csibm278", "csibm280", "csibm281",
		"csibm284", "csibm285", "csibm290", "csibm297", "csibm420", "csibm423", "csibm424",
		"csibm500", "csibm803", "csibm851", "csibm855", "csibm856", "csibm857", "csibm860",
		"csibm863", "csibm864", "csibm865", "csibm866", "csibm868", "csibm869", "csibm870",
		"csibm871", "csibm880", "csibm891", "csibm901", "csibm902", "csibm903", "csibm904",
		"csibm905", "csibm918", "csibm921", "csibm922", "csibm930", "csibm932", "csibm933",
		"csibm935", "csibm937", "csibm939", "csibm943", "csibm1008", "csibm1025", "csibm1026",
		"csibm1097", "csibm1112", "csibm1122", "csibm1123", "csibm1124", "csibm1129", "csibm1130",
		"csibm1132", "csibm1133", "csibm1137", "csibm1140", "csibm1141", "csibm1142", "csibm1143",
		"csibm1144", "csibm1145", "csibm1146", "csibm1147", "csibm1148", "csibm1149", "csibm1153",
		"csibm1154", "csibm1155", "csibm1156", "csibm1157", "csibm1158", "csibm1160", "csibm1161",
		"csibm1163", "csibm1164", "csibm1166", "csibm1167", "csibm1364", "csibm1371", "csibm1388",
		"csibm1390", "csibm1399", "csibm4517", "csibm4899", "csibm4909", "csibm4971", "csibm5347",
		"csibm9030", "csibm9066", "csibm9448", "csibm12712", "csibm16804", "csibm11621162",
		"csiso4unitedkingdom", "csiso10swedish", "csiso11swedishfornames",
		"csiso14jisc6220ro", "csiso15italian", "csiso16portugese", "csiso17spanish",
		"csiso18greek7old", "csiso19latingreek", "csiso21german", "csiso25french",
		"csiso27latingreek1", "csiso49inis", "csiso50inis8", "csiso51iniscyrillic",
		"csiso58gb1988", "csiso60danishnorwegian", "csiso60norwegian1", "csiso61norwegian2",
		"csiso69french", "csiso84portuguese2", "csiso85spanish2", "csiso86hungarian",
		"csiso88greek7", "csiso89asmo449", "csiso90", "csiso92jisc62991984b", "csiso99naplps",
		"csiso103t618bit", "csiso111ecmacyrillic", "csiso121canadian1", "csiso122canadian2",
		"csiso139csn369103", "csiso141jusib1002", "csiso143iecp271", "csiso150",
		"csiso150greekccitt", "csiso151cuba", "csiso153gost1976874", "csiso646danish",
		"csiso2022cn", "csiso2022jp", "csiso2022jp2", "csiso2022kr", "csiso2033",
		"csiso5427cyrillic", "csiso5427cyrillic1981", "csiso5428greek", "csiso10367box",
		"csisolatin1", "csisolatin2", "csisolatin3", "csisolatin4", "csisolatin5", "csisolatin6",
		"csisolatinarabic", "csisolatincyrillic", "csisolatingreek", "csisolatinhebrew",
		"cskoi8r", "csksc5636", "csmacintosh", "csnatsdano", "csnatssefi", "csn_369103",
		"cspc8codepage437", "cspc775baltic", "cspc850multilingual", "cspc862latinhebrew",
		"cspcp852", "csshiftjis", "csucs4", "csunicode", "cswindows31j", "cuba", "cwi-2", "cwi",
		"cyrillic", "de", "dec-mcs", "dec", "decmcs", "din_66003", "dk", "ds2089", "ds_2089", "e13b",
		"ebcdic-at-de-a", "ebcdic-at-de", "ebcdic-be", "ebcdic-br", "ebcdic-ca-fr",
		"ebcdic-cp-ar1", "ebcdic-cp-ar2", "ebcdic-cp-be", "ebcdic-cp-ca", "ebcdic-cp-ch",
		"ebcdic-cp-dk", "ebcdic-cp-es", "ebcdic-cp-fi", "ebcdic-cp-fr", "ebcdic-cp-gb",
		"ebcdic-cp-gr", "ebcdic-cp-he", "ebcdic-cp-is", "ebcdic-cp-it", "ebcdic-cp-nl",
		"ebcdic-cp-no", "ebcdic-cp-roece", "ebcdic-cp-se", "ebcdic-cp-tr", "ebcdic-cp-us",
		"ebcdic-cp-wt", "ebcdic-cp-yu", "ebcdic-cyrillic", "ebcdic-dk-no-a", "ebcdic-dk-no",
		"ebcdic-es-a", "ebcdic-es-s", "ebcdic-es", "ebcdic-fi-se-a", "ebcdic-fi-se", "ebcdic-fr",
		"ebcdic-greek", "ebcdic-int", "ebcdic-int1", "ebcdic-is-friss", "ebcdic-it",
		"ebcdic-jp-e", "ebcdic-jp-kana", "ebcdic-pt", "ebcdic-uk", "ebcdic-us", "ebcdicatde",
		"ebcdicatdea", "ebcdiccafr", "ebcdicdkno", "ebcdicdknoa", "ebcdices", "ebcdicesa",
		"ebcdicess", "ebcdicfise", "ebcdicfisea", "ebcdicfr", "ebcdicisfriss", "ebcdicit",
		"ebcdicpt", "ebcdicuk", "ebcdicus", "ecma-114", "ecma-118", "ecma-128", "ecma-cyrillic",
		"ecmacyrillic", "elot_928", "es", "es2", "euc-cn", "euc-jisx0213", "euc-jp-ms", "euc-jp",
		"euc-kr", "euc-tw", "euccn", "eucjp-ms", "eucjp-open", "eucjp-win", "eucjp", "euckr", "euctw",
		"fi", "fr", "gb", "gb2312", "gb13000", "gb18030", "gbk", "gb_1988-80", "gb_198880",
		"georgian-academy", "georgian-ps", "gost_19768-74", "gost_19768", "gost_1976874",
		"greek-ccitt", "greek", "greek7-old", "greek7", "greek7old", "greek8", "greekccitt",
		"hebrew", "hp-greek8", "hp-roman8", "hp-roman9", "hp-thai8", "hp-turkish8", "hpgreek8",
		"hproman8", "hproman9", "hpthai8", "hpturkish8", "hu", "ibm-803", "ibm-856", "ibm-901",
		"ibm-902", "ibm-921", "ibm-922", "ibm-930", "ibm-932", "ibm-933", "ibm-935", "ibm-937",
		"ibm-939", "ibm-943", "ibm-1008", "ibm-1025", "ibm-1046", "ibm-1047", "ibm-1097", "ibm-1112",
		"ibm-1122", "ibm-1123", "ibm-1124", "ibm-1129", "ibm-1130", "ibm-1132", "ibm-1133",
		"ibm-1137", "ibm-1140", "ibm-1141", "ibm-1142", "ibm-1143", "ibm-1144", "ibm-1145",
		"ibm-1146", "ibm-1147", "ibm-1148", "ibm-1149", "ibm-1153", "ibm-1154", "ibm-1155",
		"ibm-1156", "ibm-1157", "ibm-1158", "ibm-1160", "ibm-1161", "ibm-1162", "ibm-1163",
		"ibm-1164", "ibm-1166", "ibm-1167", "ibm-1364", "ibm-1371", "ibm-1388", "ibm-1390",
		"ibm-1399", "ibm-4517", "ibm-4899", "ibm-4909", "ibm-4971", "ibm-5347", "ibm-9030",
		"ibm-9066", "ibm-9448", "ibm-12712", "ibm-16804", "ibm037", "ibm038", "ibm256", "ibm273",
		"ibm274", "ibm275", "ibm277", "ibm278", "ibm280", "ibm281", "ibm284", "ibm285", "ibm290",
		"ibm297", "ibm367", "ibm420", "ibm423", "ibm424", "ibm437", "ibm500", "ibm775", "ibm803",
		"ibm813", "ibm819", "ibm848", "ibm850", "ibm851", "ibm852", "ibm855", "ibm856", "ibm857",
		"ibm860", "ibm861", "ibm862", "ibm863", "ibm864", "ibm865", "ibm866", "ibm866nav", "ibm868",
		"ibm869", "ibm870", "ibm871", "ibm874", "ibm875", "ibm880", "ibm891", "ibm901", "ibm902",
		"ibm903", "ibm904", "ibm905", "ibm912", "ibm915", "ibm916", "ibm918", "ibm920", "ibm921",
		"ibm922", "ibm930", "ibm932", "ibm933", "ibm935", "ibm937", "ibm939", "ibm943", "ibm1004",
		"ibm1008", "ibm1025", "ibm1026", "ibm1046", "ibm1047", "ibm1089", "ibm1097", "ibm1112",
		"ibm1122", "ibm1123", "ibm1124", "ibm1129", "ibm1130", "ibm1132", "ibm1133", "ibm1137",
		"ibm1140", "ibm1141", "ibm1142", "ibm1143", "ibm1144", "ibm1145", "ibm1146", "ibm1147",
		"ibm1148", "ibm1149", "ibm1153", "ibm1154", "ibm1155", "ibm1156", "ibm1157", "ibm1158",
		"ibm1160", "ibm1161", "ibm1162", "ibm1163", "ibm1164", "ibm1166", "ibm1167", "ibm1364",
		"ibm1371", "ibm1388", "ibm1390", "ibm1399", "ibm4517", "ibm4899", "ibm4909", "ibm4971",
		"ibm5347", "ibm9030", "ibm9066", "ibm9448", "ibm12712", "ibm16804", "iec_p27-1", "iec_p271",
		"inis-8", "inis-cyrillic", "inis", "inis8", "iniscyrillic", "isiri-3342", "isiri3342",
		"iso-2022-cn-ext", "iso-2022-cn", "iso-2022-jp-2", "iso-2022-jp-3", "iso-2022-jp",
		"iso-2022-kr", "iso-8859-1", "iso-8859-2", "iso-8859-3", "iso-8859-4", "iso-8859-5",
		"iso-8859-6", "iso-8859-7", "iso-8859-8", "iso-8859-9", "iso-8859-9e", "iso-8859-10",
		"iso-8859-11", "iso-8859-13", "iso-8859-14", "iso-8859-15", "iso-8859-16", "iso-10646",
		"iso-10646/ucs2", "iso-10646/ucs4", "iso-10646/utf-8", "iso-10646/utf8", "iso-celtic",
		"iso-ir-4", "iso-ir-6", "iso-ir-8-1", "iso-ir-9-1", "iso-ir-10", "iso-ir-11", "iso-ir-14",
		"iso-ir-15", "iso-ir-16", "iso-ir-17", "iso-ir-18", "iso-ir-19", "iso-ir-21", "iso-ir-25",
		"iso-ir-27", "iso-ir-37", "iso-ir-49", "iso-ir-50", "iso-ir-51", "iso-ir-54", "iso-ir-55",
		"iso-ir-57", "iso-ir-60", "iso-ir-61", "iso-ir-69", "iso-ir-84", "iso-ir-85", "iso-ir-86",
		"iso-ir-88", "iso-ir-89", "iso-ir-90", "iso-ir-92", "iso-ir-98", "iso-ir-99", "iso-ir-100",
		"iso-ir-101", "iso-ir-103", "iso-ir-109", "iso-ir-110", "iso-ir-111", "iso-ir-121",
		"iso-ir-122", "iso-ir-126", "iso-ir-127", "iso-ir-138", "iso-ir-139", "iso-ir-141",
		"iso-ir-143", "iso-ir-144", "iso-ir-148", "iso-ir-150", "iso-ir-151", "iso-ir-153",
		"iso-ir-155", "iso-ir-156", "iso-ir-157", "iso-ir-166", "iso-ir-179", "iso-ir-193",
		"iso-ir-197", "iso-ir-199", "iso-ir-203", "iso-ir-209", "iso-ir-226", "iso/tr_11548-1",
		"iso646-ca", "iso646-ca2", "iso646-cn", "iso646-cu", "iso646-de", "iso646-dk", "iso646-es",
		"iso646-es2", "iso646-fi", "iso646-fr", "iso646-fr1", "iso646-gb", "iso646-hu",
		"iso646-it", "iso646-jp-ocr-b", "iso646-jp", "iso646-kr", "iso646-no", "iso646-no2",
		"iso646-pt", "iso646-pt2", "iso646-se", "iso646-se2", "iso646-us", "iso646-yu",
		"iso2022cn", "iso2022cnext", "iso2022jp", "iso2022jp2", "iso2022kr", "iso6937",
		"iso8859-1", "iso8859-2", "iso8859-3", "iso8859-4", "iso8859-5", "iso8859-6", "iso8859-7",
		"iso8859-8", "iso8859-9", "iso8859-9e", "iso8859-10", "iso8859-11", "iso8859-13",
		"iso8859-14", "iso8859-15", "iso8859-16", "iso11548-1", "iso88591", "iso88592", "iso88593",
		"iso88594", "iso88595", "iso88596", "iso88597", "iso88598", "iso88599", "iso88599e",
		"iso885910", "iso885911", "iso885913", "iso885914", "iso885915", "iso885916",
		"iso_646.irv:1991", "iso_2033-1983", "iso_2033", "iso_5427-ext", "iso_5427",
		"iso_5427:1981", "iso_5427ext", "iso_5428", "iso_5428:1980", "iso_6937-2",
		"iso_6937-2:1983", "iso_6937", "iso_6937:1992", "iso_8859-1", "iso_8859-1:1987",
		"iso_8859-2", "iso_8859-2:1987", "iso_8859-3", "iso_8859-3:1988", "iso_8859-4",
		"iso_8859-4:1988", "iso_8859-5", "iso_8859-5:1988", "iso_8859-6", "iso_8859-6:1987",
		"iso_8859-7", "iso_8859-7:1987", "iso_8859-7:2003", "iso_8859-8", "iso_8859-8:1988",
		"iso_8859-9", "iso_8859-9:1989", "iso_8859-9e", "iso_8859-10", "iso_8859-10:1992",
		"iso_8859-14", "iso_8859-14:1998", "iso_8859-15", "iso_8859-15:1998", "iso_8859-16",
		"iso_8859-16:2001", "iso_9036", "iso_10367-box", "iso_10367box", "iso_11548-1",
		"iso_69372", "it", "jis_c6220-1969-ro", "jis_c6229-1984-b", "jis_c62201969ro",
		"jis_c62291984b", "johab", "jp-ocr-b", "jp", "js", "jus_i.b1.002", "koi-7", "koi-8", "koi8-r",
		"koi8-ru", "koi8-t", "koi8-u", "koi8", "koi8r", "koi8u", "ksc5636", "l1", "l2", "l3", "l4", "l5", "l6",
		"l7", "l8", "l10", "latin-9", "latin-greek-1", "latin-greek", "latin1", "latin2", "latin3",
		"latin4", "latin5", "latin6", "latin7", "latin8", "latin9", "latin10", "latingreek",
		"latingreek1", "mac-centraleurope", "mac-cyrillic", "mac-is", "mac-sami", "mac-uk", "mac",
		"maccyrillic", "macintosh", "macis", "macuk", "macukrainian", "mik", "ms-ansi", "ms-arab",
		"ms-cyrl", "ms-ee", "ms-greek", "ms-hebr", "ms-mac-cyrillic", "ms-turk", "ms932", "ms936",
		"mscp949", "mscp1361", "msmaccyrillic", "msz_7795.3", "ms_kanji", "naplps", "nats-dano",
		"nats-sefi", "natsdano", "natssefi", "nc_nc0010", "nc_nc00-10", "nc_nc00-10:81",
		"nf_z_62-010", "nf_z_62-010_(1973)", "nf_z_62-010_1973", "nf_z_62010",
		"nf_z_62010_1973", "no", "no2", "ns_4551-1", "ns_4551-2", "ns_45511", "ns_45512",
		"os2latin1", "osf00010001", "osf00010002", "osf00010003", "osf00010004", "osf00010005",
		"osf00010006", "osf00010007", "osf00010008", "osf00010009", "osf0001000a", "osf00010020",
		"osf00010100", "osf00010101", "osf00010102", "osf00010104", "osf00010105", "osf00010106",
		"osf00030010", "osf0004000a", "osf0005000a", "osf05010001", "osf100201a4", "osf100201a8",
		"osf100201b5", "osf100201f4", "osf100203b5", "osf1002011c", "osf1002011d", "osf1002035d",
		"osf1002035e", "osf1002035f", "osf1002036b", "osf1002037b", "osf10010001", "osf10010004",
		"osf10010006", "osf10020025", "osf10020111", "osf10020115", "osf10020116", "osf10020118",
		"osf10020122", "osf10020129", "osf10020352", "osf10020354", "osf10020357", "osf10020359",
		"osf10020360", "osf10020364", "osf10020365", "osf10020366", "osf10020367", "osf10020370",
		"osf10020387", "osf10020388", "osf10020396", "osf10020402", "osf10020417", "pt", "pt2",
		"pt154", "r8", "r9", "rk1048", "roman8", "roman9", "ruscii", "se", "se2", "sen_850200_b",
		"sen_850200_c", "shift-jis", "shift_jis", "shift_jisx0213", "sjis-open", "sjis-win",
		"sjis", "ss636127", "strk1048-2002", "st_sev_358-88", "t.61-8bit", "t.61", "t.618bit",
		"tcvn-5712", "tcvn", "tcvn5712-1", "tcvn5712-1:1993", "thai8", "tis-620", "tis620-0",
		"tis620.2529-1", "tis620.2533-0", "tis620", "ts-5881", "tscii", "turkish8", "ucs-2",
		"ucs-2be", "ucs-2le", "ucs-4", "ucs-4be", "ucs-4le", "ucs2", "ucs4", "uhc", "ujis", "uk",
		"unicode", "unicodebig", "unicodelittle", "us-ascii", "us", "utf-7", "utf-8", "utf-16",
		"utf-16be", "utf-16le", "utf-32", "utf-32be", "utf-32le", "utf7", "utf8", "utf16", "utf16be",
		"utf16le", "utf32", "utf32be", "utf32le", "viscii", "wchar_t", "win-sami-2", "winbaltrim",
		"windows-31j", "windows-874", "windows-936", "windows-1250", "windows-1251",
		"windows-1252", "windows-1253", "windows-1254", "windows-1255", "windows-1256",
		"windows-1257", "windows-1258", "winsami2", "ws2", "yu"}
}
