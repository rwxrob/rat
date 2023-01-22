package x

var (
	SyntaxError = `"%!ERROR: invalid rat/x type or syntax"`
	UsageName   = `"%!USAGE: x.Name{name, rule}"`
	UsageSave   = `"%!USAGE: x.Save{name}"`
	UsageVal    = `"%!USAGE: x.Val{name}"`
	UsageRef    = `"%!USAGE: x.Ref{name}"`
	UsageIs     = `"%!USAGE: namedFunc or x.IsFunc or x.Is{namedFunc}"`
	UsageAny    = `"%!USAGE: x.Any{n} or x.Any{m, n} or x.Any{m, 0}"`
	UsageStr    = `"%!USAGE: x.Str{...any}"`
	UsageSeq    = `"%!USAGE: x.Seq{...rule}"`
	UsageOne    = `"%!USAGE: x.One{...rule}"`
	UsageOpt    = `"%!USAGE: x.Opt{rule}"`
	UsageMn1    = `"%!USAGE: x.Mn1{rule}"`
	UsageMn0    = `"%!USAGE: x.Mn0{rule}"`
	UsageMin    = `"%!USAGE: x.Min{n, rule}"`
	UsageMax    = `"%!USAGE: x.Max{n, rule}"`
	UsageMmx    = `"%!USAGE: x.Mmx{n, m, rule}"`
	UsageRep    = `"%!USAGE: x.Rep{n, rule}"`
	UsageSee    = `"%!USAGE: x.See{rule}"`
	UsageNot    = `"%!USAGE: x.Not{rule}"`
	UsageTo     = `"%!USAGE: x.To{rule}"`
	UsageRng    = `"%!USAGE: x.Rng{beg, end}"`
	UsageEnd    = `"%!USAGE: x.End{}"`
)
