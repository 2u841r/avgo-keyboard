package main

type BengaliChar struct {
	Bengali string
	IsVowel bool
}

type KeyMap struct {
	Patterns        map[string]BengaliChar
	VowelDiacritics map[string]string
}

func NewKeyMap() *KeyMap {
	patterns := make(map[string]BengaliChar)
	vowelDiacritics := make(map[string]string)

	// Independent vowels (স্বরবর্ণ)
	patterns["o"] = BengaliChar{Bengali: "অ", IsVowel: true}
	patterns["a"] = BengaliChar{Bengali: "আ", IsVowel: true}
	patterns["i"] = BengaliChar{Bengali: "ই", IsVowel: true}
	patterns["I"] = BengaliChar{Bengali: "ঈ", IsVowel: true}
	patterns["u"] = BengaliChar{Bengali: "উ", IsVowel: true}
	patterns["U"] = BengaliChar{Bengali: "ঊ", IsVowel: true}
	patterns["rri"] = BengaliChar{Bengali: "ঋ", IsVowel: true}
	patterns["e"] = BengaliChar{Bengali: "এ", IsVowel: true}
	patterns["oi"] = BengaliChar{Bengali: "ঐ", IsVowel: true}
	patterns["O"] = BengaliChar{Bengali: "ও", IsVowel: true}
	patterns["ou"] = BengaliChar{Bengali: "ঔ", IsVowel: true}

	// Vowel diacritics (কার) - used after consonants
	vowelDiacritics["o"] = ""    // অ-কার (inherent vowel - no diacritic needed)
	vowelDiacritics["a"] = "া"   // আ-কার
	vowelDiacritics["i"] = "ি"   // ই-কার
	vowelDiacritics["I"] = "ী"   // ঈ-কার
	vowelDiacritics["u"] = "ু"   // উ-কার
	vowelDiacritics["U"] = "ূ"   // ঊ-কার
	vowelDiacritics["rri"] = "ৃ" // ঋ-কার
	vowelDiacritics["e"] = "ে"   // এ-কার
	vowelDiacritics["oi"] = "ৈ"  // ঐ-কার
	vowelDiacritics["O"] = "ো"   // ও-কার
	vowelDiacritics["ou"] = "ৌ"  // ঔ-কার

	// র-ফলা (r-phola) patterns - consonant + r
	patterns["phr"] = BengaliChar{Bengali: "ফ্র", IsVowel: false}
	patterns["bhr"] = BengaliChar{Bengali: "ভ্র", IsVowel: false}
	patterns["thr"] = BengaliChar{Bengali: "থ্র", IsVowel: false}
	patterns["dhr"] = BengaliChar{Bengali: "ধ্র", IsVowel: false}
	patterns["shr"] = BengaliChar{Bengali: "শ্র", IsVowel: false}
	patterns["chr"] = BengaliChar{Bengali: "ছ্র", IsVowel: false}
	patterns["pr"] = BengaliChar{Bengali: "প্র", IsVowel: false}
	patterns["br"] = BengaliChar{Bengali: "ব্র", IsVowel: false}
	patterns["tr"] = BengaliChar{Bengali: "ত্র", IsVowel: false}
	patterns["dr"] = BengaliChar{Bengali: "দ্র", IsVowel: false}
	patterns["kr"] = BengaliChar{Bengali: "ক্র", IsVowel: false}
	patterns["gr"] = BengaliChar{Bengali: "গ্র", IsVowel: false}
	patterns["jr"] = BengaliChar{Bengali: "জ্র", IsVowel: false}
	patterns["mr"] = BengaliChar{Bengali: "ম্র", IsVowel: false}
	patterns["nr"] = BengaliChar{Bengali: "ন্র", IsVowel: false}
	patterns["sr"] = BengaliChar{Bengali: "স্র", IsVowel: false}
	patterns["hr"] = BengaliChar{Bengali: "হ্র", IsVowel: false}
	patterns["fr"] = BengaliChar{Bengali: "ফ্র", IsVowel: false}
	patterns["vr"] = BengaliChar{Bengali: "ভ্র", IsVowel: false}
	patterns["lr"] = BengaliChar{Bengali: "ল্র", IsVowel: false}
	patterns["rr"] = BengaliChar{Bengali: "র্", IsVowel: false}
	patterns["Tr"] = BengaliChar{Bengali: "ট্র", IsVowel: false}
	patterns["Dr"] = BengaliChar{Bengali: "ড্র", IsVowel: false}
	patterns["Nr"] = BengaliChar{Bengali: "ণ্র", IsVowel: false}

	// Complex letters (conjuncts) - MUST come before simple consonants
	patterns["shk"] = BengaliChar{Bengali: "ষ্ক", IsVowel: false}
	patterns["shkr"] = BengaliChar{Bengali: "ষ্ক্র", IsVowel: false}
	patterns["kSh"] = BengaliChar{Bengali: "ক্ষ", IsVowel: false}
	patterns["kkh"] = BengaliChar{Bengali: "ক্ষ", IsVowel: false}
	patterns["jY"] = BengaliChar{Bengali: "জ্ঞ", IsVowel: false}
	patterns["gg"] = BengaliChar{Bengali: "জ্ঞ", IsVowel: false}

	// Double consonants
	patterns["kk"] = BengaliChar{Bengali: "ক্ক", IsVowel: false}
	patterns["kT"] = BengaliChar{Bengali: "ক্ট", IsVowel: false}
	patterns["kt"] = BengaliChar{Bengali: "ক্ত", IsVowel: false}
	patterns["kw"] = BengaliChar{Bengali: "ক্ব", IsVowel: false}
	patterns["km"] = BengaliChar{Bengali: "ক্ম", IsVowel: false}
	patterns["kl"] = BengaliChar{Bengali: "ক্ল", IsVowel: false}
	patterns["ks"] = BengaliChar{Bengali: "ক্স", IsVowel: false}

	patterns["tt"] = BengaliChar{Bengali: "ত্ত", IsVowel: false}
	patterns["tn"] = BengaliChar{Bengali: "ত্ন", IsVowel: false}
	patterns["tw"] = BengaliChar{Bengali: "ত্ব", IsVowel: false}
	patterns["tm"] = BengaliChar{Bengali: "ত্ম", IsVowel: false}

	patterns["dd"] = BengaliChar{Bengali: "দ্দ", IsVowel: false}
	patterns["dw"] = BengaliChar{Bengali: "দ্ব", IsVowel: false}
	patterns["dm"] = BengaliChar{Bengali: "দ্ম", IsVowel: false}

	patterns["nn"] = BengaliChar{Bengali: "ন্ন", IsVowel: false}
	patterns["nt"] = BengaliChar{Bengali: "ন্ত", IsVowel: false}
	patterns["nd"] = BengaliChar{Bengali: "ন্দ", IsVowel: false}
	patterns["nw"] = BengaliChar{Bengali: "ন্ব", IsVowel: false}
	patterns["nm"] = BengaliChar{Bengali: "ন্ম", IsVowel: false}

	patterns["pp"] = BengaliChar{Bengali: "প্প", IsVowel: false}
	patterns["pt"] = BengaliChar{Bengali: "প্ত", IsVowel: false}
	patterns["pl"] = BengaliChar{Bengali: "প্ল", IsVowel: false}

	patterns["bb"] = BengaliChar{Bengali: "ব্ব", IsVowel: false}
	patterns["bd"] = BengaliChar{Bengali: "ব্দ", IsVowel: false}
	patterns["bl"] = BengaliChar{Bengali: "ব্ল", IsVowel: false}

	patterns["mm"] = BengaliChar{Bengali: "ম্ম", IsVowel: false}
	patterns["mp"] = BengaliChar{Bengali: "ম্প", IsVowel: false}
	patterns["mb"] = BengaliChar{Bengali: "ম্ব", IsVowel: false}
	patterns["ml"] = BengaliChar{Bengali: "ম্ল", IsVowel: false}

	patterns["ll"] = BengaliChar{Bengali: "ল্ল", IsVowel: false}
	patterns["lk"] = BengaliChar{Bengali: "ল্ক", IsVowel: false}
	patterns["lg"] = BengaliChar{Bengali: "ল্গ", IsVowel: false}
	patterns["lp"] = BengaliChar{Bengali: "ল্প", IsVowel: false}
	patterns["lw"] = BengaliChar{Bengali: "ল্ব", IsVowel: false}
	patterns["lm"] = BengaliChar{Bengali: "ল্ম", IsVowel: false}

	patterns["sk"] = BengaliChar{Bengali: "স্ক", IsVowel: false}
	patterns["st"] = BengaliChar{Bengali: "স্ত", IsVowel: false}
	patterns["sn"] = BengaliChar{Bengali: "স্ন", IsVowel: false}
	patterns["sp"] = BengaliChar{Bengali: "স্প", IsVowel: false}
	patterns["sw"] = BengaliChar{Bengali: "স্ব", IsVowel: false}
	patterns["sm"] = BengaliChar{Bengali: "স্ম", IsVowel: false}
	patterns["sl"] = BengaliChar{Bengali: "স্ল", IsVowel: false}

	// Consonants (ব্যঞ্জনবর্ণ) - Order matters: longer patterns first
	patterns["kh"] = BengaliChar{Bengali: "খ", IsVowel: false}
	patterns["gh"] = BengaliChar{Bengali: "ঘ", IsVowel: false}
	patterns["ch"] = BengaliChar{Bengali: "ছ", IsVowel: false}
	patterns["jh"] = BengaliChar{Bengali: "ঝ", IsVowel: false}
	patterns["Th"] = BengaliChar{Bengali: "ঠ", IsVowel: false}
	patterns["Dh"] = BengaliChar{Bengali: "ঢ", IsVowel: false}
	patterns["th"] = BengaliChar{Bengali: "থ", IsVowel: false}
	patterns["dh"] = BengaliChar{Bengali: "ধ", IsVowel: false}
	patterns["ph"] = BengaliChar{Bengali: "ফ", IsVowel: false}
	patterns["bh"] = BengaliChar{Bengali: "ভ", IsVowel: false}
	patterns["Rh"] = BengaliChar{Bengali: "ঢ়", IsVowel: false}
	patterns["ya"] = BengaliChar{Bengali: "য়া", IsVowel: false}
	patterns["Ng"] = BengaliChar{Bengali: "ঙ", IsVowel: false}
	patterns["ng"] = BengaliChar{Bengali: "ং", IsVowel: false}
	patterns[".t"] = BengaliChar{Bengali: "ৎ", IsVowel: false}
	patterns[".n"] = BengaliChar{Bengali: "ঁ", IsVowel: false}

	// Single consonants
	patterns["k"] = BengaliChar{Bengali: "ক", IsVowel: false}
	patterns["g"] = BengaliChar{Bengali: "গ", IsVowel: false}
	patterns["C"] = BengaliChar{Bengali: "ছ", IsVowel: false}
	patterns["c"] = BengaliChar{Bengali: "চ", IsVowel: false}
	patterns["j"] = BengaliChar{Bengali: "জ", IsVowel: false}
	patterns["Y"] = BengaliChar{Bengali: "ঞ", IsVowel: false}
	patterns["T"] = BengaliChar{Bengali: "ট", IsVowel: false}
	patterns["D"] = BengaliChar{Bengali: "ড", IsVowel: false}
	patterns["N"] = BengaliChar{Bengali: "ণ", IsVowel: false}
	patterns["t"] = BengaliChar{Bengali: "ত", IsVowel: false}
	patterns["d"] = BengaliChar{Bengali: "দ", IsVowel: false}
	patterns["n"] = BengaliChar{Bengali: "ন", IsVowel: false}
	patterns["f"] = BengaliChar{Bengali: "ফ", IsVowel: false}
	patterns["p"] = BengaliChar{Bengali: "প", IsVowel: false}
	patterns["v"] = BengaliChar{Bengali: "ভ", IsVowel: false}
	patterns["b"] = BengaliChar{Bengali: "ব", IsVowel: false}
	patterns["m"] = BengaliChar{Bengali: "ম", IsVowel: false}
	patterns["z"] = BengaliChar{Bengali: "য", IsVowel: false}
	patterns["r"] = BengaliChar{Bengali: "র", IsVowel: false}
	patterns["l"] = BengaliChar{Bengali: "ল", IsVowel: false}
	patterns["Sh"] = BengaliChar{Bengali: "ষ", IsVowel: false}
	patterns["sh"] = BengaliChar{Bengali: "ষ", IsVowel: false}
	patterns["S"] = BengaliChar{Bengali: "শ", IsVowel: false}
	patterns["s"] = BengaliChar{Bengali: "স", IsVowel: false}
	patterns["h"] = BengaliChar{Bengali: "হ", IsVowel: false}
	patterns["R"] = BengaliChar{Bengali: "ড়", IsVowel: false}
	patterns["y"] = BengaliChar{Bengali: "য়", IsVowel: false}
	patterns["yo"] = BengaliChar{Bengali: "য়", IsVowel: false}
	patterns[":"] = BengaliChar{Bengali: "ঃ", IsVowel: false}
	patterns["H"] = BengaliChar{Bengali: "ঃ", IsVowel: false}

	// Numbers and others
	patterns["0"] = BengaliChar{Bengali: "০", IsVowel: false}
	patterns["1"] = BengaliChar{Bengali: "১", IsVowel: false}
	patterns["2"] = BengaliChar{Bengali: "২", IsVowel: false}
	patterns["3"] = BengaliChar{Bengali: "৩", IsVowel: false}
	patterns["4"] = BengaliChar{Bengali: "৪", IsVowel: false}
	patterns["5"] = BengaliChar{Bengali: "৫", IsVowel: false}
	patterns["6"] = BengaliChar{Bengali: "৬", IsVowel: false}
	patterns["7"] = BengaliChar{Bengali: "৭", IsVowel: false}
	patterns["8"] = BengaliChar{Bengali: "৮", IsVowel: false}
	patterns["9"] = BengaliChar{Bengali: "৯", IsVowel: false}
	patterns["."] = BengaliChar{Bengali: "।", IsVowel: false}
	patterns["$"] = BengaliChar{Bengali: "৳", IsVowel: false}
	patterns["aya"] = BengaliChar{Bengali: "অ্যা", IsVowel: false}

	return &KeyMap{
		Patterns:        patterns,
		VowelDiacritics: vowelDiacritics,
	}
}
