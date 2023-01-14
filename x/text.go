package x

var (
	SyntaxError = `"%!ERROR: invalid rat/x type or syntax"`
	UsageName   = `"%!USAGE: x.Name{name, rule}"`
	UsageID     = `"%!USAGE: x.ID{id, rule}"`
	UsageRef    = `"%!USAGE: x.Ref{name}"`
	UsageRid    = `"%!USAGE: x.Rid{id}"`
	UsageIs     = `"%!USAGE: namedFunc or x.IsFunc or x.Is{namedFunc}"`
	UsageAny    = `"%!USAGE: x.Any{n} or x.Any{n, m}"`
	UsageLit    = `"%!USAGE: x.Lit{...any}"`
	UsageSeq    = `"%!USAGE: x.Seq{...rule}"`
	UsageOne    = `"%!USAGE: x.One{...rule}"`
	UsageOpt    = `"%!USAGE: x.Opt{rule}"`
	UsageMn1    = `"%!USAGE: x.Mn1{rule}"`
	UsageMn0    = `"%!USAGE: x.Mn0{rule}"`
	UsageMin    = `"%!USAGE: x.Min{n, rule}"`
	UsageMax    = `"%!USAGE: x.Max{n, rule}"`
	UsageMmx    = `"%!USAGE: x.Mmx{n, m, rule}"`
	UsageRep    = `"%!USAGE: x.Rep{n, rule}"`
	UsagePos    = `"%!USAGE: x.Pos{rule}"`
	UsageNeg    = `"%!USAGE: x.Neg{rule}"`
	UsageToi    = `"%!USAGE: x.Toi{rule}"`
	UsageTox    = `"%!USAGE: x.Tox{rule}"`
	UsageRng    = `"%!USAGE: x.Rng{beg, end}"`
	UsageEnd    = `"%!USAGE: x.End{}"`
)
