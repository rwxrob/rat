package x

const (
	_SyntaxError = `"%!ERROR: invalid rat/x type or syntax"`
	_UsageRule   = `"%!USAGE: x.Rule{name, intid, rule}"`
	_UsageRef    = `"%!USAGE: x.Ref{name}"`
	_UsageIs     = `"%!USAGE: x.Is{namedfunc}"`
	_UsageAny    = `"%!USAGE: x.Any{n}"`
	_UsageLit    = `"%!USAGE: x.Lit{...any}"`
	_UsageSeq    = `"%!USAGE: x.Seq{...rule}"`
	_UsageOne    = `"%!USAGE: x.One{...rule}"`
	_UsageOpt    = `"%!USAGE: x.Opt{rule}"`
	_UsageMn1    = `"%!USAGE: x.Mn1{rule}"`
	_UsageMn0    = `"%!USAGE: x.Mn0{rule}"`
	_UsageMin    = `"%!USAGE: x.Min{n, rule}"`
	_UsageMax    = `"%!USAGE: x.Max{n, rule}"`
	_UsageMmx    = `"%!USAGE: x.Mmx{n, m, rule}"`
	_UsageRep    = `"%!USAGE: x.Rep{n, rule}"`
	_UsagePos    = `"%!USAGE: x.Pos{rule}"`
	_UsageNeg    = `"%!USAGE: x.Neg{rule}"`
	_UsageToi    = `"%!USAGE: x.Toi{rule}"`
	_UsageTox    = `"%!USAGE: x.Tox{rule}"`
	_UsageRng    = `"%!USAGE: x.Rng{beg, end}"`
	_UsageEnd    = `"%!USAGE: x.End{}"`
)
