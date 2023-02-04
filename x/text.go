package x

var (
	SyntaxError = `"%!ERROR: invalid rat/x type or syntax"`
	UsageN      = `"%!USAGE: x.N{name, rule}"`
	UsageSav    = `"%!USAGE: x.Sav{name}"`
	UsageVal    = `"%!USAGE: x.Val{name}"`
	UsageRef    = `"%!USAGE: x.Ref{name}"`
	UsageIs     = `"%!USAGE: namedFunc or x.IsFunc or x.Is{namedFunc}"`
	UsageAny    = `"%!USAGE: x.Any{n} or x.Any{m, n} or x.Any{m, 0}"`
	UsageStr    = `"%!USAGE: x.Str{...any}"`
	UsageSeq    = `"%!USAGE: x.Seq{...rule}"`
	UsageOne    = `"%!USAGE: x.One{...rule}"`
	UsageMmx    = `"%!USAGE: x.Mmx{m, n, rule}"`
	UsageSee    = `"%!USAGE: x.See{rule}"`
	UsageNot    = `"%!USAGE: x.Not{rule}"`
	UsageTo     = `"%!USAGE: x.To{rule}"`
	UsageRng    = `"%!USAGE: x.Rng{beg, end}"`
	UsageEnd    = `"%!USAGE: x.End{}"`
)
