package btrfs

//	offsets := [...]uint64{
//		0,
//		4194304,
//		12582912,
//		20971520,
//		29360128,
//		1103101952,
//		2176843776,
//		3250585600,
//		4324327424,
//		5398069248,
//		6471811072,
//		7545552896,
//		8619294720,
//		9693036544,
//		10766778368,
//		11840520192,
//		12914262016,
//		13988003840,
//		15061745664,
//		16135487488,
//		17209229312,
//		18282971136,
//		19356712960,
//		20430454784,
//		21504196608,
//		22577938432,
//		23651680256,
//		24725422080,
//		25799163904,
//		26872905728,
//		27946647552,
//		29020389376,
//		30094131200,
//		32241614848,
//		33315356672,
//		34389098496,
//		35462840320,
//		36536582144,
//		37610323968,
//		38684065792,
//		39757807616,
//		40831549440,
//		41905291264,
//		42979033088,
//		44052774912,
//		45126516736,
//		47274000384,
//		48347742208,
//		49421484032,
//		50495225856,
//		52642709504,
//		53716451328,
//		54790193152,
//		55863934976,
//		58011418624,
//		61232644096,
//		62306385920,
//		65527611392,
//		66601353216,
//		67675095040,
//		68748836864,
//		71970062336,
//		73043804160,
//		74117545984,
//		75191287808,
//		76265029632,
//		77338771456,
//		79486255104,
//		80559996928,
//		83781222400,
//		84854964224,
//		85928706048,
//		87002447872,
//		88076189696,
//		89149931520,
//		90223673344,
//		91297415168,
//		92371156992,
//		93444898816,
//		94518640640,
//		95592382464,
//		98813607936,
//		105256058880,
//		106329800704,
//		107403542528,
//		111698509824,
//		112772251648,
//		114919735296,
//		115993477120,
//		118140960768,
//		119214702592,
//		120288444416,
//		123509669888,
//		124583411712,
//		125657153536,
//		126730895360,
//		128878379008,
//		129952120832,
//		131025862656,
//		132099604480,
//		135320829952,
//		137468313600,
//		138542055424,
//		141763280896,
//		142837022720,
//		143910764544,
//		144984506368,
//		146058248192,
//		148205731840,
//		149279473664,
//		150353215488,
//		151426957312,
//		152500699136,
//		153574440960,
//		154648182784,
//		156795666432,
//		158943150080,
//		160016891904,
//		161090633728,
//		162164375552,
//		163238117376,
//		164311859200,
//		165385601024,
//		166459342848,
//		167533084672,
//		168606826496,
//		169680568320,
//		170754310144,
//		171828051968,
//		172901793792,
//		173975535616,
//		175049277440,
//		176123019264,
//		178270502912,
//		179344244736,
//		180417986560,
//		181491728384,
//		182565470208,
//		183639212032,
//		184712953856,
//		185786695680,
//		186860437504,
//		187934179328,
//		189007921152,
//		191155404800,
//		192229146624,
//		193302888448,
//		195450372096,
//		196524113920,
//		198671597568,
//		199745339392,
//		200819081216,
//		201892823040,
//		204040306688,
//		205114048512,
//		206187790336,
//		207261532160,
//		208335273984,
//		209409015808,
//		210482757632,
//		211556499456,
//		212630241280,
//		213703983104,
//		214777724928,
//		215851466752,
//		216925208576,
//		217998950400,
//		219072692224,
//		220146434048,
//		221220175872,
//		223367659520,
//		224441401344,
//		225515143168,
//		226588884992,
//		227662626816,
//		228736368640,
//		229810110464,
//		230883852288,
//		231957594112,
//		233031335936,
//		234105077760,
//		235178819584,
//		236252561408,
//		237326303232,
//		238400045056,
//		239473786880,
//		240547528704,
//		241621270528,
//		242695012352,
//		243768754176,
//		244842496000,
//		245916237824,
//		246989979648,
//		248063721472,
//		249137463296,
//		250211205120,
//		251284946944,
//		252358688768,
//		253432430592,
//		254506172416,
//		255579914240,
//		256653656064,
//		257727397888,
//		258801139712,
//		259874881536,
//		260948623360,
//		262022365184,
//		263096107008,
//		264169848832,
//		265243590656,
//		266317332480,
//		267391074304,
//		268464816128,
//		269538557952,
//		270612299776,
//		271686041600,
//		272759783424,
//		273833525248,
//		274907267072,
//		275981008896,
//		277054750720,
//		278128492544,
//		279202234368,
//		280275976192,
//		281349718016,
//		282423459840,
//		283497201664,
//		284570943488,
//		285644685312,
//		286718427136,
//		287792168960,
//		288865910784,
//		289939652608,
//		291013394432,
//		292087136256,
//		293160878080,
//		294234619904,
//		295308361728,
//		296382103552,
//		297455845376,
//		298529587200,
//		299603329024,
//		300677070848,
//		301750812672,
//		302824554496,
//		303898296320,
//		304972038144,
//		306045779968,
//		307119521792,
//		308193263616,
//		309267005440,
//		310340747264,
//		311414489088,
//		312488230912,
//		313561972736,
//		314635714560,
//		315709456384,
//		316783198208,
//		317856940032,
//		318930681856,
//		320004423680,
//		321078165504,
//		322151907328,
//		323225649152,
//		324299390976,
//		325373132800,
//		326446874624,
//		327520616448,
//		328594358272,
//		329668100096,
//		330741841920,
//		331815583744,
//		332889325568,
//		333963067392,
//		335036809216,
//		336110551040,
//		337184292864,
//		338258034688,
//		339331776512,
//		340405518336,
//		341479260160,
//		342553001984,
//		343626743808,
//		344700485632,
//		345774227456,
//		346847969280,
//		347921711104,
//		348995452928,
//		350069194752,
//		351142936576,
//		352216678400,
//		353290420224,
//		354364162048,
//		355437903872,
//		356511645696,
//		357585387520,
//		358659129344,
//		359732871168,
//		360806612992,
//		361880354816,
//		362954096640,
//		364027838464,
//		365101580288,
//		366175322112,
//		367249063936,
//		368322805760,
//		369396547584,
//		370470289408,
//		371544031232,
//		372617773056,
//		373691514880,
//		374765256704,
//		375838998528,
//		376912740352,
//		377986482176,
//		379060224000,
//		380133965824,
//		381207707648,
//		382281449472,
//		383355191296,
//		384428933120,
//		385502674944,
//		386576416768,
//		387650158592,
//		388723900416,
//		389797642240,
//		390871384064,
//		391945125888,
//		393018867712,
//		394092609536,
//		395166351360,
//		396240093184,
//		397313835008,
//		398387576832,
//		399461318656,
//		400535060480,
//		401608802304,
//		402682544128,
//		403756285952,
//		404830027776,
//		405903769600,
//		406977511424,
//		408051253248,
//		409124995072,
//		410198736896,
//		411272478720,
//		412346220544,
//		413419962368,
//		414493704192,
//		415567446016,
//		416641187840,
//		417714929664,
//		418788671488,
//		419862413312,
//		420936155136,
//		422009896960,
//		423083638784,
//		424157380608,
//		425231122432,
//		426304864256,
//		427378606080,
//		428452347904,
//		429526089728,
//		430599831552,
//		431673573376,
//		432747315200,
//		433821057024,
//		434894798848,
//		435968540672,
//		437042282496,
//		438116024320,
//		439189766144,
//		440263507968,
//		441337249792,
//		442410991616,
//		443484733440,
//		444558475264,
//		445632217088,
//		446705958912,
//		447779700736,
//		448853442560,
//		449927184384,
//		451000926208,
//		452074668032,
//		453148409856,
//		454222151680,
//		455295893504,
//		456369635328,
//		457443377152,
//		458517118976,
//		459590860800,
//		460664602624,
//		461738344448,
//		462812086272,
//		463885828096,
//		464959569920,
//		466033311744,
//		467107053568,
//		468180795392,
//		469254537216,
//		470328279040,
//		471402020864,
//		472475762688,
//		473549504512,
//		474623246336,
//		475696988160,
//		476770729984,
//		477844471808,
//		478918213632,
//		479991955456,
//		481065697280,
//		482139439104,
//		483213180928,
//		484286922752,
//		485360664576,
//		486434406400,
//		487508148224,
//		488581890048,
//		489655631872,
//		490729373696,
//		491803115520,
//		492876857344,
//		493950599168,
//		495024340992,
//		496098082816,
//		497171824640,
//		498245566464,
//		499319308288,
//		500393050112,
//		501466791936,
//		502540533760,
//		503614275584,
//		504688017408,
//		505761759232,
//		506835501056,
//		507909242880,
//		508982984704,
//		510056726528,
//		511130468352,
//		512204210176,
//		513277952000,
//		514351693824,
//		515425435648,
//		516499177472,
//		517572919296,
//		518646661120,
//		519720402944,
//		520794144768,
//		521867886592,
//		522941628416,
//		524015370240,
//		525089112064,
//		526162853888,
//		527236595712,
//		528310337536,
//		529384079360,
//		530457821184,
//		531531563008,
//		532605304832,
//		533679046656,
//		534752788480,
//		535826530304,
//		536900272128,
//		537974013952,
//		539047755776,
//		540121497600,
//		541195239424,
//		542268981248,
//		543342723072,
//		544416464896,
//		545490206720,
//		546563948544,
//		547637690368,
//		548711432192,
//		549785174016,
//		550858915840,
//		551932657664,
//		553006399488,
//		554080141312,
//		555153883136,
//		556227624960,
//		557301366784,
//		558375108608,
//		559448850432,
//		560522592256,
//		561596334080,
//		562670075904,
//		563743817728,
//		564817559552,
//		565891301376,
//		566965043200,
//		568038785024,
//		569112526848,
//		570186268672,
//		571260010496,
//		572333752320,
//		573407494144,
//		574481235968,
//		575554977792,
//		576628719616,
//		577702461440,
//		578776203264,
//		579849945088,
//		580923686912,
//		581997428736,
//		583071170560,
//		584144912384,
//		585218654208,
//		586292396032,
//		587366137856,
//		588439879680,
//		589513621504,
//		590587363328,
//		591661105152,
//		592734846976,
//		593808588800,
//		594882330624,
//		595956072448,
//		597029814272,
//		598103556096,
//		599177297920,
//		600251039744,
//		601324781568,
//		602398523392,
//		603472265216,
//		604546007040,
//		605619748864,
//		606693490688,
//		607767232512,
//		608840974336,
//		609914716160,
//		610988457984,
//		612062199808,
//		613135941632,
//		614209683456,
//		615283425280,
//		616357167104,
//		617430908928,
//		618504650752,
//		619578392576,
//		620652134400,
//		621725876224,
//		622799618048,
//		623873359872,
//		624947101696,
//		626020843520,
//		627094585344,
//		628168327168,
//		629242068992,
//		630315810816,
//		631389552640,
//		632463294464,
//		633537036288,
//		634610778112,
//		635684519936,
//		636758261760,
//		637832003584,
//		638905745408,
//		639979487232,
//		641053229056,
//		642126970880,
//		643200712704,
//		644274454528,
//		645348196352,
//		646421938176,
//		647495680000,
//		648569421824,
//		649643163648,
//		650716905472,
//		651790647296,
//		652864389120,
//		653938130944,
//		655011872768,
//		656085614592,
//		657159356416,
//		658233098240,
//		659306840064,
//		660380581888,
//		661454323712,
//		662528065536,
//		663601807360,
//		664675549184,
//		665749291008,
//		666823032832,
//		667896774656,
//		668970516480,
//		670044258304,
//		671118000128,
//		672191741952,
//		673265483776,
//		674339225600,
//		675412967424,
//		676486709248,
//		677560451072,
//		678634192896,
//		679707934720,
//		680781676544,
//		681855418368,
//		682929160192,
//		684002902016,
//		685076643840,
//		686150385664,
//		687224127488,
//		688297869312,
//		689371611136,
//		690445352960,
//		691519094784,
//		692592836608,
//		693666578432,
//		694740320256,
//		695814062080,
//		696887803904,
//		697961545728,
//		699035287552,
//		700109029376,
//		701182771200,
//		702256513024,
//		703330254848,
//		704403996672,
//		705477738496,
//		706551480320,
//		707625222144,
//		708698963968,
//		709772705792,
//		710846447616,
//		711920189440,
//		712993931264,
//		714067673088,
//		715141414912,
//		716215156736,
//		717288898560,
//		718362640384,
//		719436382208,
//		720510124032,
//		721583865856,
//		722657607680,"encoding/hex"
//		723731349504,
//		724805091328,
//		725878833152,
//		726952574976,
//		728026316800,
//		729100058624,
//		730173800448,
//		731247542272,
//		732321284096,
//		733395025920,
//		734468767744,
//		735542509568,
//		736616251392,
//		737689993216,
//		738763735040,
//		739837476864,
//		740911218688,
//		741984960512,
//		743058702336,
//		744132444160,
//		745206185984,
//		746279927808,
//		747353669632,
//		748427411456,
//		749501153280,
//		750574895104,
//		751648636928,
//		752722378752,
//		753796120576,
//		754869862400,
//		755943604224,
//		757017346048,
//		758091087872,
//		759164829696,
//		760238571520,
//		761312313344,
//		762386055168,
//		763459796992,
//		764533538816,
//		765607280640,
//		766681022464,
//		767754764288,
//		768828506112,
//		769902247936,
//		770975989760,
//		772049731584,
//		773123473408,
//		774197215232,
//		775270957056,
//		776344698880,
//		777418440704,
//		778492182528,
//		779565924352,
//		780639666176,
//		781713408000,
//		782787149824,
//		783860891648,
//		784934633472,
//		786008375296,
//		787082117120,
//		788155858944,
//		789229600768,
//		790303342592,
//		791377084416,
//		792450826240,
//		793524568064,
//		794598309888,
//		795672051712,
//		796745793536,
//		797819535360,
//		798893277184,
//		799967019008,
//		801040760832,
//		802114502656,
//		803188244480,
//		804261986304,
//		805335728128,
//		806409469952,
//		807483211776,
//		808556953600,
//		809630695424,
//		810704437248,
//		811778179072,
//		812851920896,
//		813925662720,
//		814999404544,
//		816073146368,
//		817146888192,
//		818220630016,
//		819294371840,
//		820368113664,
//		821441855488,
//		822515597312,
//		823589339136,
//		824663080960,
//		825736822784,
//		826810564608,
//		827884306432,
//		828958048256,
//		830031790080,
//		831105531904,
//		832179273728,
//		833253015552,
//		834326757376,
//		835400499200,
//		836474241024,
//		837547982848,
//		838621724672,
//		839695466496,
//		840769208320,
//		841842950144,
//		842916691968,
//		843990433792,
//		845064175616,
//		846137917440,
//		847211659264,
//		848285401088,
//		849359142912,
//		850432884736,
//		851506626560,
//		852580368384,
//		853654110208,
//		854727852032,
//		855801593856,
//		856875335680,
//		857949077504,
//		859022819328,
//		860096561152,
//		861170302976,
//		862244044800,
//		863317786624,
//		864391528448,
//		865465270272,
//		866539012096,
//		867612753920,
//		868686495744,
//		869760237568,
//		870833979392,
//		871907721216,
//		872981463040,
//		874055204864,
//		875128946688,
//		876202688512,
//		877276430336,
//		878350172160,
//		879423913984,
//		880497655808,
//		881571397632,
//		882645139456,
//		883718881280,
//		884792623104,
//		885866364928,
//		886940106752,
//		888013848576,
//		889087590400,
//		890161332224,
//		891235074048,
//		892308815872,
//		893382557696,
//		894456299520,
//		895530041344,
//		896603783168,
//		897677524992,
//		898751266816,
//		899825008640,
//		900898750464,
//		901972492288,
//		903046234112,
//		904119975936,
//		905193717760,
//		906267459584,
//		907341201408,
//		908414943232,
//		909488685056,
//		910562426880,
//		911636168704,
//		912709910528,
//		913783652352,
//		914857394176,
//		915931136000,
//		917004877824,
//		918078619648,
//		919152361472,
//		920226103296,
//		921299845120,
//		922373586944,
//		923447328768,
//		924521070592,
//		925594812416,
//		926668554240,
//		927742296064,
//		928816037888,
//		929889779712,
//		930963521536,
//		932037263360,
//		933111005184,
//		934184747008,
//		935258488832,
//		936332230656,
//		937405972480,
//		938479714304,
//		939553456128,
//		940627197952,
//		941700939776,
//		942774681600,
//		943848423424,
//		944922165248,
//		945995907072,
//		947069648896,
//		948143390720,
//		949217132544,
//		950290874368,
//		951364616192,
//		952438358016,
//		953512099840,
//		954585841664,
//		955659583488,
//		956733325312,
//		957807067136,
//		958880808960,
//		959954550784,
//		961028292608,
//		962102034432,
//		963175776256,
//		964249518080,
//		965323259904,
//		966397001728,
//		967470743552,
//		968544485376,
//		969618227200,
//		970691969024,
//		971765710848,
//		972839452672,
//		973913194496,
//		974986936320,
//	}
//	var empty []byte = make([]byte, btrfs.BTRFS_SIZE)
// print some info while scanning for metadata of the fsid
//func processBlock(bytenr uint64, header *btrfs.BtrfsHeader, rc *btrfs.RecoverControl, fsid []byte) {
//
//	if !btrfs.IsSuperBlockAddress(uint64(bytenr)) {
//		if btrfs.BtrfsReadHeader(rc.Fd, header, bytenr) {
//			if bytes.Equal(header.Fsid[:], fsid) {
//				fmt.Printf("Btrfs header @%v:\n%+v\n\n", bytenr, header)
//				decodeObjectID(int64(header.Owner))
//
//				if ret, items := btrfs.BtrfsReadItems(rc.Fd, header.Nritems, bytenr); ret {
//					for _, item := range items {
//						decodeKeyID(item.Key.Type)
//					}
//				}
//			}
//		}
//	}
//}

//func decodeObjectID(id int64) {
//	// Object IDS
//	/* holds pointers to all of the tree roots */
//	switch id {
//	case btrfs.BTRFS_ROOT_TREE_OBJECTID:
//		fmt.Println("BTRFS_ROOT_TREE_OBJECTID")
//
//	/* stores information about which extents are in use, and reference counts */
//	case btrfs.BTRFS_EXTENT_TREE_OBJECTID:
//		fmt.Println("BTRFS_EXTENT_TREE_OBJECTID")
//
//	/*
//	 * chunk tree stores translations from logical -> physical block numbering
//	 * the super block points to the chunk tree
//	 */
//	case btrfs.BTRFS_CHUNK_TREE_OBJECTID:
//		fmt.Println("BTRFS_CHUNK_TREE_OBJECTID")
//
//	/*
//	 * stores information about which areas of a given device are in use.
//	 * one per device.  The tree of tree roots points to the device tree
//	 */
//	case btrfs.BTRFS_DEV_TREE_OBJECTID:
//		fmt.Println("BTRFS_DEV_TREE_OBJECTID")
//
//	/* one per subvolume, storing files and directories */
//	case btrfs.BTRFS_FS_TREE_OBJECTID:
//		fmt.Println("BTRFS_FS_TREE_OBJECTID")
//
//	/* directory objectid inside the root tree */
//	case btrfs.BTRFS_ROOT_TREE_DIR_OBJECTID:
//		fmt.Println("BTRFS_ROOT_TREE_DIR_OBJECTID")
//	/* holds checksums of all the data extents */
//	case btrfs.BTRFS_CSUM_TREE_OBJECTID:
//		fmt.Println("BTRFS_CSUM_TREE_OBJECTID")
//	case btrfs.BTRFS_QUOTA_TREE_OBJECTID:
//		fmt.Println("BTRFS_QUOTA_TREE_OBJECTID")
//
//	/* for storing items that use the BTRFS_UUID_KEY* */
//	case btrfs.BTRFS_UUID_TREE_OBJECTID:
//		fmt.Println("BTRFS_UUID_TREE_OBJECTID")
//
//	/* for storing balance parameters in the root tree */
//	case btrfs.BTRFS_BALANCE_OBJECTID:
//		fmt.Println("BTRFS_BALANCE_OBJECTID")
//
//	/* oprhan objectid for tracking unlinked/truncated files */
//	case btrfs.BTRFS_ORPHAN_OBJECTID:
//		fmt.Println("BTRFS_ORPHAN_OBJECTID")
//
//	/* does write ahead logging to speed up fsyncs */
//	case btrfs.BTRFS_TREE_LOG_OBJECTID:
//		fmt.Println("BTRFS_TREE_LOG_OBJECTID")
//	case btrfs.BTRFS_TREE_LOG_FIXUP_OBJECTID:
//		fmt.Println("BTRFS_TREE_LOG_FIXUP_OBJECTID")
//
//	/* space balancing */
//	case btrfs.BTRFS_TREE_RELOC_OBJECTID:
//		fmt.Println("BTRFS_TREE_RELOC_OBJECTID")
//	case btrfs.BTRFS_DATA_RELOC_TREE_OBJECTID:
//		fmt.Println("BTRFS_DATA_RELOC_TREE_OBJECTID")
//
//	/*
//	 * extent checksums all have this objectid
//	 * this allows them to share the logging tree
//	 * for fsyncs
//	 */
//	case btrfs.BTRFS_EXTENT_CSUM_OBJECTID:
//		fmt.Println("BTRFS_EXTENT_CSUM_OBJECTID")
//
//	/* For storing free space cache */
//	case btrfs.BTRFS_FREE_SPACE_OBJECTID:
//		fmt.Println("BTRFS_FREE_SPACE_OBJECTID")
//
//	/*
//	 * The inode number assigned to the special inode for sotring
//	 * free ino cache
//	 */
//	case btrfs.BTRFS_FREE_INO_OBJECTID:
//		fmt.Println("BTRFS_FREE_INO_OBJECTID")
//
//	/* dummy objectid represents multiple objectids */
//	case btrfs.BTRFS_MULTIPLE_OBJECTIDS:
//		fmt.Println("BTRFS_MULTIPLE_OBJECTIDS")
//
//	/*
//	 * All files have objectids in this range.
//	 */
//	//	case btrfs.BTRFS_FIRST_FREE_OBJECTID:
//	//		fmt.Println("BTRFS_FIRST_FREE_OBJECTID")
//	//	case btrfs.BTRFS_LAST_FREE_OBJECTID:
//	//		fmt.Println("BTRFS_LAST_FREE_OBJECTID")
//	//	case btrfs.BTRFS_FIRST_CHUNK_TREE_OBJECTID:
//	//		fmt.Println("BTRFS_FIRST_CHUNK_TREE_OBJECTID")
//	//
//	//	/*
//	//	 * the device items go into the chunk tree.  The key is in the form
//	//	 * [ 1 BTRFS_DEV_ITEM_KEY deviceId ]
//	//	 */
//	//	case btrfs.BTRFS_DEV_ITEMS_OBJECTID:
//	//		fmt.Println("BTRFS_DEV_ITEMS_OBJECTID")
//	default:
//		if (id > 0 && id > btrfs.BTRFS_FIRST_FREE_OBJECTID) || (id < 0 && id < btrfs.BTRFS_LAST_FREE_OBJECTID) {
//			fmt.Printf("Numbered Object %08x\n", uint64(id))
//		} else {
//			fmt.Println("UNKNOWN OBJECTID")
//
//		}
//	}
//}
//
//func decodeKeyID(id uint8) {
//	/*
//	 * inode items have the data typically returned from stat and store other
//	 * info about object characteristics.  There is one for every file and dir in
//	 * the FS
//	 */
//	switch id {
//	case btrfs.BTRFS_INODE_ITEM_KEY:
//		fmt.Println("BTRFS_INODE_ITEM_KEY")
//	case btrfs.BTRFS_INODE_REF_KEY:
//		fmt.Println("BTRFS_INODE_REF_KEY")
//	case btrfs.BTRFS_INODE_EXTREF_KEY:
//		fmt.Println("BTRFS_INODE_EXTREF_KEY")
//	case btrfs.BTRFS_XATTR_ITEM_KEY:
//		fmt.Println("BTRFS_XATTR_ITEM_KEY")
//	case btrfs.BTRFS_ORPHAN_ITEM_KEY:
//		fmt.Println("BTRFS_ORPHAN_ITEM_KEY")
//
//	case btrfs.BTRFS_DIR_LOG_ITEM_KEY:
//		fmt.Println("BTRFS_DIR_LOG_ITEM_KEY")
//	case btrfs.BTRFS_DIR_LOG_INDEX_KEY:
//		fmt.Println("BTRFS_DIR_LOG_INDEX_KEY")
//	/*
//	 * dir items are the name -> inode pointers in a directory.  There is one
//	 * for every name in a directory.
//	 */
//	case btrfs.BTRFS_DIR_ITEM_KEY:
//		fmt.Println("BTRFS_DIR_ITEM_KEY")
//	case btrfs.BTRFS_DIR_INDEX_KEY:
//		fmt.Println("BTRFS_DIR_INDEX_KEY")
//
//	/*
//	 * extent data is for file data
//	 */
//	case btrfs.BTRFS_EXTENT_DATA_KEY:
//		fmt.Println("BTRFS_EXTENT_DATA_KEY")
//
//	/*
//	 * csum items have the checksums for data in the extents
//	 */
//	case btrfs.BTRFS_CSUM_ITEM_KEY:
//		fmt.Println("BTRFS_CSUM_ITEM_KEY")
//	/*
//	 * extent csums are stored in a separate tree and hold csums for
//	 * an entire extent on disk.
//	 */
//	case btrfs.BTRFS_EXTENT_CSUM_KEY:
//		fmt.Println("BTRFS_EXTENT_CSUM_KEY")
//
//	/*
//	 * root items point to tree roots.  There are typically in the root
//	 * tree used by the super block to find all the other trees
//	 */
//	case btrfs.BTRFS_ROOT_ITEM_KEY:
//		fmt.Println("BTRFS_ROOT_ITEM_KEY")
//
//	/*
//	 * root backrefs tie subvols and snapshots to the directory entries that
//	 * reference them
//	 */
//	case btrfs.BTRFS_ROOT_BACKREF_KEY:
//		fmt.Println("BTRFS_ROOT_BACKREF_KEY")
//
//	/*
//	 * root refs make a fast index for listing all of the snapshots and
//	 * subvolumes referenced by a given root.  They point directly to the
//	 * directory item in the root that references the subvol
//	 */
//	case btrfs.BTRFS_ROOT_REF_KEY:
//		fmt.Println("BTRFS_ROOT_REF_KEY")
//
//	/*
//	 * extent items are in the extent map tree.  These record which blocks
//	 * are used, and how many references there are to each block
//	 */
//	case btrfs.BTRFS_EXTENT_ITEM_KEY:
//		fmt.Println("BTRFS_EXTENT_ITEM_KEY")
//
//	/*
//	 * The same as the BTRFS_EXTENT_ITEM_KEY, except it's metadata we already know
//	 * the length, so we save the level in key->offset instead of the length.
//	 */
//	case btrfs.BTRFS_METADATA_ITEM_KEY:
//		fmt.Println("BTRFS_METADATA_ITEM_KEY")
//
//	case btrfs.BTRFS_TREE_BLOCK_REF_KEY:
//		fmt.Println("BTRFS_TREE_BLOCK_REF_KEY")
//
//	case btrfs.BTRFS_EXTENT_DATA_REF_KEY:
//		fmt.Println("BTRFS_EXTENT_DATA_REF_KEY")
//
//	/* old style extent backrefs */
//	case btrfs.BTRFS_EXTENT_REF_V0_KEY:
//		fmt.Println("BTRFS_EXTENT_REF_V0_KEY")
//
//	case btrfs.BTRFS_SHARED_BLOCK_REF_KEY:
//		fmt.Println("BTRFS_SHARED_BLOCK_REF_KEY")
//
//	case btrfs.BTRFS_SHARED_DATA_REF_KEY:
//		fmt.Println("BTRFS_SHARED_DATA_REF_KEY")
//
//	/*
//	 * block groups give us hints into the extent allocation trees.  Which
//	 * blocks are free etc etc
//	 */
//	case btrfs.BTRFS_BLOCK_GROUP_ITEM_KEY:
//		fmt.Println("BTRFS_BLOCK_GROUP_ITEM_KEY")
//
//	case btrfs.BTRFS_DEV_EXTENT_KEY:
//		fmt.Println("BTRFS_DEV_EXTENT_KEY")
//	case btrfs.BTRFS_DEV_ITEM_KEY:
//		fmt.Println("BTRFS_DEV_ITEM_KEY")
//	case btrfs.BTRFS_CHUNK_ITEM_KEY:
//		fmt.Println("BTRFS_CHUNK_ITEM_KEY")
//
//	case btrfs.BTRFS_BALANCE_ITEM_KEY:
//		fmt.Println("BTRFS_BALANCE_ITEM_KEY")
//
//	/*
//	 * quota groups
//	 */
//	case btrfs.BTRFS_QGROUP_STATUS_KEY:
//		fmt.Println("BTRFS_QGROUP_STATUS_KEY")
//	case btrfs.BTRFS_QGROUP_INFO_KEY:
//		fmt.Println("BTRFS_QGROUP_INFO_KEY")
//	case btrfs.BTRFS_QGROUP_LIMIT_KEY:
//		fmt.Println("BTRFS_QGROUP_LIMIT_KEY")
//	case btrfs.BTRFS_QGROUP_RELATION_KEY:
//		fmt.Println("BTRFS_QGROUP_RELATION_KEY")
//
//	/*
//	 * Persistently stores the io stats in the device tree.
//	 * One key for all stats, (0, BTRFS_DEV_STATS_KEY, devid).
//	 */
//	case btrfs.BTRFS_DEV_STATS_KEY:
//		fmt.Println("BTRFS_DEV_STATS_KEY")
//
//	/*
//	 * Persistently stores the device replace state in the device tree.
//	 * The key is built like this: (0, BTRFS_DEV_REPLACE_KEY, 0).
//	 */
//	case btrfs.BTRFS_DEV_REPLACE_KEY:
//		fmt.Println("BTRFS_DEV_REPLACE_KEY")
//
//		/*
//		 * Stores items that allow to quickly map UUIDs to something else.
//		 * These items are part of the filesystem UUID tree.
//		 * The key is built like this:
//		 * (UUIDUpper_64Bits, BTRFS_UUID_KEY*, UUIDLower_64Bits).
//		 */
//		//#if BTRFS_UUID_SIZE != 16
//		//#error "UUID items require case btrfs.BTRFS_UUID_SIZE:
//		fmt.Println("BTRFS_UUID_SIZE")
//	//#endif
//	case btrfs.BTRFS_UUID_KEY_SUBVOL:
//		fmt.Println("BTRFS_UUID_KEY_SUBVOL")
//	case btrfs.BTRFS_UUID_KEY_RECEIVED_SUBVOL:
//		fmt.Println("BTRFS_UUID_KEY_RECEIVED_SUBVOL")
//		/* received subvols */
//
//	/*
//	 * string items are for debugging.  They just store a short string of
//	 * data in the FS
//	 */
//	case btrfs.BTRFS_STRING_ITEM_KEY:
//		fmt.Println("BTRFS_STRING_ITEM_KEY")
//		//	default:
//		//		fmt.Println("UNKNOWN ITEM KEY ID")
//
//	}
//}
//
//func extractMetadataRecord(rc *RecoverControl, leaf *ExtentBuffer) bool {
//	var (
//		//	struct btrfsKey key;
//		ret = false
//	//	int i;
//	//	u32 nritems;
//	//
//	)
//	nritems := uint64((*BtrfsHeader)(unsafe.Pointer(&leaf.Data)).Nritems)
//	for i := uint64(0); i < nritems; i++ {
//		//		btrfsItemKeyToCpu(leaf, &key, i);
//		//		switch (key.type) {
//		//		case BTRFS_BLOCK_GROUP_ITEM_KEY:
//		//			pthreadMutexLock(&rc->rcLock);
//		//			ret = processBlockGroupItem(&rc->bg, leaf, &key, i);
//		//			pthreadMutexUnlock(&rc->rcLock);
//		//			break;
//		//		case BTRFS_CHUNK_ITEM_KEY:
//		//			pthreadMutexLock(&rc->rcLock);
//		//			ret = processChunkItem(&rc->chunk, leaf, &key, i);
//		//			pthreadMutexUnlock(&rc->rcLock);
//		//			break;
//		//		case BTRFS_DEV_EXTENT_KEY:
//		//			pthreadMutexLock(&rc->rcLock);
//		//			ret = processDeviceExtentItem(&rc->devext, leaf,
//		//							 &key, i);
//		//			pthreadMutexUnlock(&rc->rcLock);
//		//			break;
//		//		}
//		//		if (ret)
//		//			break;
//	}
//	return ret
//}
//func scanOneDevice(devScan *DeviceScan) bool {
//	var (
//		buf    *ExtentBuffer
//		bytenr uint64
//		ret    bool = false
//		//struct deviceScan *devScan = (struct deviceScan *)devScanStruct;
//		rc *RecoverControl = devScan.Rc
//		//		device *BtrfsDevice    = devScan.Dev
//		fd int = devScan.Fd
//		//		oldtype int
//		h *BtrfsHeader
//	)
//	buf = new(ExtentBuffer)
//	buf.Len = uint64(rc.Leafsize)
//	buf.Data = make([]byte, rc.Leafsize)
//loop:
//	for bytenr = 0; ; bytenr += uint64(rc.Sectorsize) {
//		ret = false
//		if IsSuperBlockAddress(bytenr) {
//			bytenr += uint64(rc.Sectorsize)
//		}
//		n, err := syscall.Pread(fd, buf.Data, int64(bytenr))
//		if err != nil {
//			log.Fatalln(os.NewSyscallError("pread64", err))
//		}
//		if n < int(rc.Leafsize) {
//			break loop
//		}
//		h = (*BtrfsHeader)(unsafe.Pointer(&buf.Data))
//		if rc.FsDevices.Fsid != h.Fsid || verifyTreeBlockCsumSilent(buf, uint16(rc.CsumSize)) {
//			continue loop
//		}
//		rc.RcLock.Lock()
//		//		ret = processExtentBuffer(&rc.EbCache, buf, device, bytenr)
//		rc.RcLock.Unlock()
//		if !ret {
//			break loop
//		}
//		if h.Level != 0 {
//			switch h.Owner {
//			case BTRFS_EXTENT_TREE_OBJECTID, BTRFS_DEV_TREE_OBJECTID:
//				/* different tree use different generation */
//				if h.Generation > rc.Generation {
//					continue loop
//				}
//				if ret = extractMetadataRecord(rc, buf); !ret {
//					break loop
//				}
//			case BTRFS_CHUNK_TREE_OBJECTID:
//				if h.Generation > rc.ChunkRootGeneration {
//					continue loop
//				}
//				if ret = extractMetadataRecord(rc, buf); !ret {
//					break loop
//				}
//			}
//		}
//	}
//
//	//	close(fd);
//	//	free(buf);
//	return ret
//}

//func processExtentBuffer(ebCache *CacheTree,
//	eb *ExtentBuffer,
//	device *BtrfsDevice, offset uint64) bool {
//	struct extentRecord *rec;
//	struct extentRecord *exist;
//	struct cacheExtent *cache;
//	int ret = 0;
//
//	rec = btrfsNewExtentRecord(eb);
//	if (!rec->cache.size)
//		goto freeOut;
//again:
//	cache = lookupCacheExtent(ebCache,
//				    rec->cache.start,
//				    rec->cache.size);
//	if (cache) {
//		exist = containerOf(cache, struct extentRecord, cache);
//
//		if (exist->generation > rec->generation)
//			goto freeOut;
//		if (exist->generation == rec->generation) {
//			if (exist->cache.start != rec->cache.start ||
//			    exist->cache.size != rec->cache.size ||
//			    memcmp(exist->csum, rec->csum, BTRFS_CSUM_SIZE)) {
//				ret = -EEXIST;
//			} else {
//				BUG_ON(exist->nmirrors >= BTRFS_MAX_MIRRORS);
//				exist->devices[exist->nmirrors] = device;
//				exist->offsets[exist->nmirrors] = offset;
//				exist->nmirrors++;
//			}
//			goto freeOut;
//		}
//		removeCacheExtent(ebCache, cache);
//		free(exist);
//		goto again;
//	}
//
//	rec->devices[0] = device;
//	rec->offsets[0] = offset;
//	rec->nmirrors++;
//	ret = insertCacheExtent(ebCache, &rec->cache);
//	BUG_ON(ret);
//out:
//	return ret;
//freeOut:
//	free(rec);
//	goto out;
//	return false
//}
// read header struct from fd at bytenr
//func BtrfsReadHeader(fd int, header *BtrfsHeader, bytenr uint64) bool {
//
//	var size = binary.Size(header)
//	var byteheader = make([]byte, size)
//	var bytebr = bytes.NewReader(byteheader)
//
//	ret, _ := syscall.Pread(fd, byteheader, int64(bytenr))
//	if ret < size {
//		return false
//	}
//	_ = binary.Read(bytebr, binary.LittleEndian, header)
//
//	return true
//
//}
//
//// read items structs from fd at bytenr
//func BtrfsReadItems(fd int, n uint32, bytenr uint64) (bool, []BtrfsItem) {
//	header := new(BtrfsHeader)
//	item := new(BtrfsItem)
//	myitems := make([]BtrfsItem, n)
//	var hsize = binary.Size(header)
//	var isize = binary.Size(item)
//	var byteitems = make([]byte, uint32(isize)*n)
//	var bytebr = bytes.NewReader(byteitems)
//
//	ret, _ := syscall.Pread(fd, byteitems, int64(bytenr+uint64(hsize)))
//	if ret < isize {
//		return false, nil
//	}
//	err := binary.Read(bytebr, binary.LittleEndian, myitems)
//	if err != nil {
//		return false, nil
//	} else {
//		return true, myitems
//	}
//
//}
