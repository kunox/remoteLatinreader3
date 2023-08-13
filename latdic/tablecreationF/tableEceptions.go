
package tablecreation 

import (
	S "latindictionary/structs"
)

// 受動態を能動態に移し、受動態をブランクに
// 能動態完了形を　PPL + SUM　に（PPLは受動感分詞を使用）
// 受動完了分詞を能動に移動。
// 未来分詞はそのまま（受動未来分詞は、受動ではなくむしろ動形容詞である）
// 
func exceptionalDEP(ret [][]string) [][]string {
	for y:=6; y < 9;y++ {		// IND Passive-> Active
		for x:=0; x<6 ; x++ {
			ret[x][y-6] = ret[x][y]
		}
	}
	for y:=13; y<15  ;y++ {		// Subj Passive -> Active
		for x:=0 ; x<6 ; x++ {
			ret[x][y-4] = ret[x][y]
			ret[x][y] = "-"
		}
	}
	for y:=17; y<19 ; y++ {		 // Imp 
		for x:=0; x <6 ; x++ {
			ret[x][y-2] = ret[x][y]
			ret[x][y] = "-"
		}
	}
	ppltemp := ret[4][20]
	ret[1][20] = ret[4][20]
	ret[4][20] = "-"
	for y:=3; y<6 ; y++ {
		ret[1][y] = ppltemp
		ret[2][y] = "+ sum"
	}
	ret[0][19] = ret[3][19]		//不定詞
	ret[3][19] = "-"
	return ret
}



func exceptionalSEM(ret [][]string) [][]string {
	for y:= 3; y<9 ; y++ {
		for x:=0; x <6 ; x++ {
			ret[x][y] = "-"
		}
	}
	for y:=11; y <15 ; y++ {
		for x:=0; x < 6 ; x++ {
			ret[x][y] = "-"
		}
	}
	for y:=17; y < 19; y++ {
		for x:=0; x < 6; x++ {
			ret[x][y] = "-"
		}
	}
	ppltemp := ret[4][20]
	ret[1][20] = ppltemp
	ret[4][20] = "";        // 受動完了分詞を能動に
	ret[5][20] = "";        // 受動未来分詞削除
	ret[1][19] = "";        // 能動完了不定詞削除
	ret[3][19] = "";        // 受動現在分詞なし
	for y:=3; y <6; y++ {
		ret[1][y] =ppltemp
		ret[2][y] = "+ sum"
	}
	for y:= 11; y<13; y++ {
		ret[1][y] = ppltemp
		ret[2][y] = "+ sum"
	}
	return ret
}


	// 完了系を現在に
	// 直説法、接続法　現在のみ、あとは全てブランク
	// Passive Perfect ParticleもActiveに移す
func exceptionalPERFDEF(ret [][]string) [][]string {
	for y:=3 ; y<6 ; y++ {
		for x:=0 ; x<6 ; x++ {
			ret[x][y-3] = ret[x][y]
			ret[x][y] = "-"
		}
	}
	for y:=11 ; y<13 ; y++ {
		for x:=0; x<6 ; x++ {
			ret[x][y-2] = ret[x][y]
			ret[x][y] = "-"
		}
	}
	ret[1][20] = ret[4][20];            // passive perfect particleをActiveに移動
	ret[4][20] = " ";
	ret[0][19] = ret[1][19];            // perfect infinitiveをpresentに移動
	ret[1][19] = " ";
	return ret;
}

func exceptionMemini(ret [][]string) [][]string {
	ret[1][15] = "memento";
	ret[4][15] = "mementote";
	return ret;
}

func exceptionalABSUM(ret [][]string) [][]string{
	ret[2][19] = "afore";
	ret[0][20] = "absens";
	return ret;
}

func exceptionFacio(ret [][]string, stem string) [][]string {
	passive := make([][]string, 3)
	passive[0] = []string { "fīo", "fīs", "fit", "fīmus", "fītis", "fīunt" }
	passive[1] = []string { "fīēbam", "fīēbas", "fīēbat", "fīēbamus", "fīēbātis", "fīēbant" }
	passive[2] = []string { "fīām", "fīes", "fīēt", "fīēmus", "fīētis", "fīent" }
	for y:=0 ; y<3 ; y++ {
		for x:=0 ; x<6 ; x++ {
			ret[x][y+6] = stem + passive[y][x]
		}
	}
	passiveSub := make([][]string, 2)
	passiveSub[0] = []string {"fīām", "fīās", "fiat", "fiamus", "fiatis", "fiant" }
	passiveSub[1] = []string { "firem", "fierēs", "fieret", "fierēmus", "fierētis", "fierent" }
	for y:=0 ; y<2 ; y++ {
		for x:=0 ; x<6 ; x++ {
			ret[x][y+13] = stem + passive[y][x]
		}
	}
	ret[1][17] = stem + "fī";
	ret[4][17] = stem + "fīte";
	ret[1][18] = stem + "fītō";
	ret[2][18] = stem + "fītō";
	ret[4][18] = stem + "fītōte";
	ret[5][18] = stem + "fīuntō";
	ret[3][19] = stem + "fierī";
	return ret;
}

	// nolo, malo, voloの例外処理
	// (1,0) (2,0) (4,0)に変則的な語を挿入 入れる語は（0,0)から判定
func exceptionV62(ret [][]string) [][]string {
	baseword := ret[0][0]
	if baseword == "volo" {
		ret[1][0] = "vīs";
		ret[2][0] = "volt,vult";
		ret[4][0] = "vultis,voltis";
	} else if baseword == "malo" {
		ret[1][0] = "māvis";
		ret[2][0] = "māvult";
		ret[4][0] = "māvultis";
	} else if baseword== "nolo" {
		ret[1][0] = "nōn vis";
		ret[2][0] = "nōn vult";
		ret[4][0] = "nōn vultis";
	}
	return ret
}


func nounexception(table [][]string, dc S.DicContent) {
	if dc.WORD == "vīs" {
		table[0][1] = "vim";
		table[0][2] = "vīs";
		table[0][3] = "vī";
		table[0][4] = "vī";
		table[1][0] = "vīrēs";
		table[1][1] = "vīrēs";
		table[1][2] = "vīrium";
		table[1][3] = "vīribus";
		table[1][4] = "vīribus";
		table[1][5] = "vīrēs";
	} else if dc.WORD == "bōs" {
		table[1][2] = "boum";
		table[1][3] = "bōbus,būbus";
		table[1][4] = "bōbus,būbus";
	} else if dc.WORD == "pĕnas" {
		table[1][2] = "penatium"
	} else if dc.WORD == "os" {		// mouthのosはoが長音、これはboneの方
		table[1][2] = "ossium";
		table[1][2] = "ossium"
	} else if dc.WORD == "imber" {
		table[1][2] = "imbrium"	
	}
	// } else if dc.WORD == "bōs" {
	// 	table[1][2] = "boum";
	// 	table[1][3] = "bōbus";
	// 	table[1][4] = "bōbus";
	// } 
}