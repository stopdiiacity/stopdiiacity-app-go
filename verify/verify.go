package verify

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	// language=JSON
	unsafeMessage = `{"exists":true}`
	// language=JSON
	unsafeOutsideWarCompanyMessage = `{"exists":true,"outside_war_company":true}`
	// language=JSON
	safeMessage = `{"exists":false}`
)

type VerifyRequest struct {
	URLs    []string `json:"urls"`
	Version int      `json:"version"`
}

type VerifyResponse struct {
	Exists bool `json:"exists"`
}

type LinksResponse struct {
	Groups []PrefixGroup `json:"groups"`
}

type PrefixGroup struct {
	Prefix   string   `json:"prefix"`
	Prefixes []string `json:"prefixes"`
}

type Response struct {
	Body       string
	StatusCode int
}

var outsideWarCompanies = [][]string{
	{
		"https://jobs.dou.ua/companies/pin-up-tech/",
	},
	{
		"https://jobs.dou.ua/companies/epam-systems/",        // https://dou.ua/forums/topic/36742/
		"https://jobs.dou.ua/companies/epam-anywhere/",       // https://dou.ua/forums/topic/36742/
		"https://djinni.co/jobs/company-epam-systems-bb0df/", // https://dou.ua/forums/topic/36742/
		"https://www.linkedin.com/company/epam-systems/",     // https://dou.ua/forums/topic/36742/
	},
	{
		"https://jobs.dou.ua/companies/andersen/",
		"https://djinni.co/jobs/company-andersen-e6bc8/",
		"https://www.linkedin.com/company/andersen-softwaredev/",
	},
	{
		"https://jobs.dou.ua/companies/dataart/",        // https://dou.ua/forums/topic/36742/
		"https://djinni.co/jobs/company-dataart-291a6/", // https://dou.ua/forums/topic/36742/
		"https://www.linkedin.com/company/dataart/",     // https://dou.ua/forums/topic/36742/
	},
	{
		"https://jobs.dou.ua/companies/nix-solutions-ltd/",    // https://dou.ua/forums/topic/36742/
		"https://djinni.co/jobs/company-nix-solutions-fe08e/", // https://dou.ua/forums/topic/36742/
		"https://www.linkedin.com/company/nix-solutions-ltd/", // https://dou.ua/forums/topic/36742/
		"https://www.linkedin.com/company/nix-community/",     // https://dou.ua/forums/topic/36742/
		"https://www.linkedin.com/company/nixs/",              // https://dou.ua/forums/topic/36742/
	},
	{
		"https://jobs.dou.ua/companies/grid-dynamics/",
		"https://djinni.co/jobs/company-grid-dynamics-0267a/",
		"https://www.linkedin.com/company/grid-dynamics/",
		"https://www.linkedin.com/company/grid-dynamics-digital-teams/",
	},
	{
		"https://jobs.dou.ua/companies/daxx-group/",
		"https://djinni.co/jobs/company-daxx-37ba9/",
	},
}

var stopDiiaCityPrefixes = []PrefixGroup{
	{
		Prefix: "https://jobs.dou.ua/companies/",
		Prefixes: []string{
			"https://jobs.dou.ua/companies/allright/",
			"https://jobs.dou.ua/companies/englishdom/",
			"https://jobs.dou.ua/companies/powercodelab/",
			"https://jobs.dou.ua/companies/powercode-academy/",
			"https://jobs.dou.ua/companies/genesis-technology-partners/", // Genesis
			"https://jobs.dou.ua/companies/solid/",                       // Genesis
			"https://jobs.dou.ua/companies/solidgate/",                   // Genesis
			"https://jobs.dou.ua/companies/headway-1/",                   // Genesis
			"https://jobs.dou.ua/companies/socialtech/",                  // Genesis
			"https://jobs.dou.ua/companies/skelar/",                      // Genesis
			"https://jobs.dou.ua/companies/obrio/",                       // Genesis
			"https://jobs.dou.ua/companies/boosters-product-team/",       // Genesis
			"https://jobs.dou.ua/companies/betterme/",                    // Genesis
			"https://jobs.dou.ua/companies/lift-stories-editor/",         // Genesis
			"https://jobs.dou.ua/companies/sendios/",                     // Genesis
			"https://jobs.dou.ua/companies/keiki/",                       // Genesis
			"https://jobs.dou.ua/companies/universe/",                    // Genesis
			"https://jobs.dou.ua/companies/amo/",                         // Genesis
			"https://jobs.dou.ua/companies/codefinity/",                  // Genesis
			"https://jobs.dou.ua/companies/intellectokids/",              // Genesis
			"https://jobs.dou.ua/companies/howly/",                       // Genesis
			"https://jobs.dou.ua/companies/appflame/",                    // Genesis
			"https://jobs.dou.ua/companies/quarks-tech/",                 // Genesis
			"https://jobs.dou.ua/companies/cpamatica/",                   // Genesis
			"https://jobs.dou.ua/companies/flyer-one-ventures/",          // Genesis
			"https://jobs.dou.ua/companies/subsub/",                      // Genesis
			"https://jobs.dou.ua/companies/plantin/",                     // Genesis
			"https://jobs.dou.ua/companies/suitsme/",                     // Genesis
			"https://jobs.dou.ua/companies/jooble/",
			"https://jobs.dou.ua/companies/netpeak/",             // netpeak
			"https://jobs.dou.ua/companies/netpeak-group/",       // netpeak
			"https://jobs.dou.ua/companies/ringostat/",           // netpeak
			"https://jobs.dou.ua/companies/serpstat/",            // netpeak
			"https://jobs.dou.ua/companies/saldo-apps/",          // netpeak
			"https://jobs.dou.ua/companies/academyocean/",        // netpeak
			"https://jobs.dou.ua/companies/tonti-laguna/",        // netpeak
			"https://jobs.dou.ua/companies/tonti-laguna-mobile/", // netpeak
			"https://jobs.dou.ua/companies/inweb/",               // netpeak
			"https://jobs.dou.ua/companies/leverx-group/",
			"https://jobs.dou.ua/companies/reface/",
			"https://jobs.dou.ua/companies/s-pro/",
			"https://jobs.dou.ua/companies/banza/",
			"https://jobs.dou.ua/companies/datagroup/",
			"https://jobs.dou.ua/companies/deus-robots/",
			"https://jobs.dou.ua/companies/parimatch-tech/", // parimatch
			"https://jobs.dou.ua/companies/gr8-tech/",       // parimatch https://dou.ua/lenta/news/parimatch-tech-turned-into-gr8-tech/
			"https://jobs.dou.ua/companies/pm-international/",
			"https://jobs.dou.ua/companies/pokermatch/",
			"https://jobs.dou.ua/companies/epam-systems/",
			"https://jobs.dou.ua/companies/epam-anywhere/",
			"https://jobs.dou.ua/companies/ajax-systems/",
			"https://jobs.dou.ua/companies/kyivstar/",
			"https://jobs.dou.ua/companies/nika-tech-family/", // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://jobs.dou.ua/companies/sigma-software/",
			"https://jobs.dou.ua/companies/sigma-technology-systems/",
			"https://jobs.dou.ua/companies/ideasoft-io/",
			"https://jobs.dou.ua/companies/intetics-co/",
			"https://jobs.dou.ua/companies/playson/",
			"https://jobs.dou.ua/companies/smartiway/",
			"https://jobs.dou.ua/companies/ciklum/",
			"https://jobs.dou.ua/companies/nix-solutions-ltd/",
			"https://jobs.dou.ua/companies/softserve/",
			"https://jobs.dou.ua/companies/raiffeisen/",
			"https://jobs.dou.ua/companies/raiffeisen-bank-international-ag/",
			"https://jobs.dou.ua/companies/eleks/",
			"https://jobs.dou.ua/companies/petcube-inc/",
			"https://jobs.dou.ua/companies/softteco/",
			"https://jobs.dou.ua/companies/luxoft/",          // https://dou.ua/lenta/articles/industry-about-diia-city/
			"https://jobs.dou.ua/companies/luxoft-training/", // https://dou.ua/lenta/articles/industry-about-diia-city/
			"https://jobs.dou.ua/companies/adtelligent/",
			"https://jobs.dou.ua/companies/intersog/",
			"https://jobs.dou.ua/companies/gismart_com/",
			"https://jobs.dou.ua/companies/softpositive/",
			"https://jobs.dou.ua/companies/vodafone-ukraine/",
			"https://jobs.dou.ua/companies/plvision/",
			"https://jobs.dou.ua/companies/treeum/",
			"https://jobs.dou.ua/companies/globallogic/",
			"https://jobs.dou.ua/companies/graintrack/",                   // inside
			"https://jobs.dou.ua/companies/intellias/",                    // question
			"https://jobs.dou.ua/companies/astound/",                      // https://dou.ua/lenta/articles/what-it-companies-think-about-bill-5376/
			"https://jobs.dou.ua/companies/wargaming/",                    // https://dou.ua/lenta/interviews/shumygora-about-wargaming/
			"https://jobs.dou.ua/companies/glovo/",                        // https://dou.ua/forums/topic/35378/
			"https://jobs.dou.ua/companies/dont-panic-recruiting-agency/", // agency, Glovo proxy https://web.archive.org/web/20211111175135/https://djinni.co/jobs/289484-frontend-engineer-v-glovo/
			"https://jobs.dou.ua/companies/staffingpartner/",              // agency, sigma & global proxy
			"https://jobs.dou.ua/companies/smart-solutions/",              // https://jobs.dou.ua/companies/smart-solutions/
			"https://jobs.dou.ua/companies/fintech-band/",                 // https://dou.ua/forums/topic/35880/
			"https://jobs.dou.ua/companies/fintech-farm/",                 // https://dou.ua/forums/topic/35880/
			"https://jobs.dou.ua/companies/itransition/",                  // https://dou.ua/forums/topic/35889/

			"https://jobs.dou.ua/companies/sendpulse/",                // https://dou.ua/lenta/articles/companies-about-diia-city/
			"https://jobs.dou.ua/companies/smart-project-gmbh/",       // https://dou.ua/lenta/articles/companies-about-diia-city/
			"https://jobs.dou.ua/companies/smart-project-llc/",        // https://dou.ua/lenta/articles/companies-about-diia-city/
			"https://jobs.dou.ua/companies/zagrava-games-by-playrix/", // https://dou.ua/lenta/articles/companies-about-diia-city/

			"https://jobs.dou.ua/companies/revolut/",       // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/
			"https://jobs.dou.ua/companies/macpaw/",        // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/
			"https://jobs.dou.ua/companies/calmerry/",      // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://jobs.dou.ua/companies/govitall/",      // calmerry https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://jobs.dou.ua/companies/innovecs/",      // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://jobs.dou.ua/companies/ilogos/",        // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://jobs.dou.ua/companies/natus-vincere/", // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469

			"https://jobs.dou.ua/companies/megogonet-/", // https://dev.ua/news/megogo-vstupaeie-v-diia-sity

			"https://jobs.dou.ua/companies/roosh/",     // https://dou.ua/lenta/news/44-companies-applied-to-join-diia-city/
			"https://jobs.dou.ua/companies/softblues/", // https://dou.ua/lenta/interviews/nuzhnyi-about-trios/

			"https://jobs.dou.ua/companies/samsung/",                                                                  // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://jobs.dou.ua/companies/ria/",                                                                      // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://jobs.dou.ua/companies/rozetka-ua-internet-supermarket/",                                          // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://jobs.dou.ua/companies/plarium/",                                                                  // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://jobs.dou.ua/companies/privatne-aktsionerne-tovaristvo-tsentr-kompyuternih-tehnologij-infoplyus/", // https://dou.ua/lenta/news/55-first-residents-of-diia-city/

			//
			//
			//
			//
			//
			//
			//
			//
			//
			//
			//

			"https://jobs.dou.ua/companies/axels/",               // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/playrix/",             // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/span-ukraine/",        // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/product-madness/",     // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/ooo-ollzap/",          // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/ipland/",              // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/vizor-games/",         // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/altexsoft/",           // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/sybenetix/",           // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/voki-gejms-ukraina/",  // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/4friends/",            // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/cqg/",                 // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/boolat-play/",         // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/ooo-dbo-soft/",        // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/druk-ua/",             // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/ooo-ppl33-35/",        // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/e-consulting/",        // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/pm-partners/",         // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/codiv-ukraine/",       // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/gk-rearden-group/",    // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/liga-zakon/",          // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/home-games/",          // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/perevaga-technology/", // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/taurus-quadra-ltd/",   // https://dou.ua/lenta/articles/diia-city-registry/
			"https://jobs.dou.ua/companies/omo-systems/",         // https://dou.ua/lenta/articles/diia-city-registry/

			"https://jobs.dou.ua/companies/trinetix/", // from email letter inside 2022-02-16
			"https://jobs.dou.ua/companies/dmarket/",  // https://jobs.dou.ua/companies/dmarket/vacancies/222376/
			"https://jobs.dou.ua/companies/comfy/",    // https://djinni.co/jobs/494285-senior-golang-developer-proektna-robota-/
			"https://jobs.dou.ua/companies/upstars/",  // https://opendatabot.ua/c/43212002
			"https://jobs.dou.ua/companies/eventyr/",  // https://jobs.dou.ua/companies/eventyr/reviews/61636/
			"https://jobs.dou.ua/companies/smarttek/", // https://opendatabot.ua/c/43961471 https://opendatabot.ua/c/40594308
			"https://jobs.dou.ua/companies/arsmoon/",  // https://opendatabot.ua/c/43961471 https://opendatabot.ua/c/40594308
			"https://jobs.dou.ua/companies/orderry/",  // https://dou.ua/lenta/articles/diia-city-registry/#2525061
		},
	},

	// https://github.com/stopdiiacity/stopdiiacity-app-go/issues/1
	// https://github.com/stopdiiacity/stopdiiacity-app-go/issues/3
	// https://github.com/stopdiiacity/stopdiiacity-app-go/issues/4
	// https://github.com/stopdiiacity/stopdiiacity-app-go/issues/5
	// https://github.com/stopdiiacity/stopdiiacity-app-go/issues/6
	// https://github.com/stopdiiacity/stopdiiacity-app-go/issues/7
	// https://github.com/stopdiiacity/stopdiiacity-app-go/issues/8
	{
		Prefix: "https://jobs.dou.ua/companies/",
		Prefixes: []string{
			"https://jobs.dou.ua/companies/reface/",
			"https://jobs.dou.ua/companies/roosh/",
			"https://jobs.dou.ua/companies/fintech-band/",
			"https://jobs.dou.ua/companies/nix-solutions-ltd/",
			"https://jobs.dou.ua/companies/softserve/",
			"https://jobs.dou.ua/companies/sigma-software/",
			"https://jobs.dou.ua/companies/ria/",
			"https://jobs.dou.ua/companies/fintech-band/",
			"https://jobs.dou.ua/companies/privatne-aktsionerne-tovaristvo-tsentr-kompyuternih-tehnologij-infoplyus/",
			"https://jobs.dou.ua/companies/axels/",
			"https://jobs.dou.ua/companies/plarium/",
			"https://jobs.dou.ua/companies/playrix/",
			"https://jobs.dou.ua/companies/span-ukraine/",
			"https://jobs.dou.ua/companies/product-madness/",
			"https://jobs.dou.ua/companies/ooo-ollzap/",
			"https://jobs.dou.ua/companies/ipland/",
			"https://jobs.dou.ua/companies/vizor-games/",
			"https://jobs.dou.ua/companies/altexsoft/",
			"https://jobs.dou.ua/companies/rozetka-ua-internet-supermarket/",
			"https://jobs.dou.ua/companies/sybenetix/",
			"https://jobs.dou.ua/companies/macpaw/",
			"https://jobs.dou.ua/companies/voki-gejms-ukraina/",
			"https://jobs.dou.ua/companies/ajax-systems/",
			"https://jobs.dou.ua/companies/4friends/",
			"https://jobs.dou.ua/companies/cqg/",
			"https://jobs.dou.ua/companies/boolat-play/",
			"https://jobs.dou.ua/companies/epam-systems/",
			"https://jobs.dou.ua/companies/samsung/",
			"https://jobs.dou.ua/companies/ooo-dbo-soft/",
			"https://jobs.dou.ua/companies/druk-ua/",
			"https://jobs.dou.ua/companies/revolut/",
			"https://jobs.dou.ua/companies/ooo-ppl33-35/",
			"https://jobs.dou.ua/companies/deus-robots/",
			"https://jobs.dou.ua/companies/e-consulting/",
			"https://jobs.dou.ua/companies/pm-partners/",
			"https://jobs.dou.ua/companies/codiv-ukraine/",
			"https://jobs.dou.ua/companies/genesis-technology-partners/",
			"https://jobs.dou.ua/companies/targer-1/",
			"https://jobs.dou.ua/companies/gk-rearden-group/",
			"https://jobs.dou.ua/companies/liga-zakon/",
			"https://jobs.dou.ua/companies/home-games/",
			"https://jobs.dou.ua/companies/perevaga-technology/",
			"https://jobs.dou.ua/companies/taurus-quadra-ltd/",
			"https://jobs.dou.ua/companies/omo-systems/",
			"https://jobs.dou.ua/companies/freysoft/",
			"https://jobs.dou.ua/companies/remed/",
			"https://jobs.dou.ua/companies/drone-ua/",
			"https://jobs.dou.ua/companies/checkbox/",
			"https://jobs.dou.ua/companies/tmsoft-ltd/",
			"https://jobs.dou.ua/companies/gameloft/",
			"https://jobs.dou.ua/companies/playgendary/",
			"https://jobs.dou.ua/companies/it-solutions/",
			"https://jobs.dou.ua/companies/zanzarra/",
			"https://jobs.dou.ua/companies/semantrum/",
			"https://jobs.dou.ua/companies/carbominer/",
			"https://jobs.dou.ua/companies/nova-poshta-tsentr/",
			"https://jobs.dou.ua/companies/tietoevry/",
			"https://jobs.dou.ua/companies/tatl-technology/",
			"https://jobs.dou.ua/companies/sixt/",
			"https://jobs.dou.ua/companies/aj-ti-artel/",
			"https://jobs.dou.ua/companies/letyshops/",
			"https://jobs.dou.ua/companies/smart-holding/",
			"https://jobs.dou.ua/companies/glomex-gmbh/",
			"https://jobs.dou.ua/companies/cinegy/",
			"https://jobs.dou.ua/companies/miratech/",
			"https://jobs.dou.ua/companies/fintech-farm/",
			"https://jobs.dou.ua/companies/brightsign/",
			"https://jobs.dou.ua/companies/demicon/",
			"https://jobs.dou.ua/companies/roosh/",
			"https://jobs.dou.ua/companies/sagax-llc/",
			"https://jobs.dou.ua/companies/it-specialist/",
			"https://jobs.dou.ua/companies/infostroy-ltd/",
			"https://jobs.dou.ua/companies/luxoft/",
			"https://jobs.dou.ua/companies/playtini/",
			"https://jobs.dou.ua/companies/tmsoft-ltd/",
			"https://jobs.dou.ua/companies/alterego/",
			"https://jobs.dou.ua/companies/rocque/",
			"https://jobs.dou.ua/companies/pandadoc/",
			"https://jobs.dou.ua/companies/cloud-services/",
			"https://jobs.dou.ua/companies/optimize-technologies/",
			"https://jobs.dou.ua/companies/megogonet-/",
			"https://jobs.dou.ua/companies/advantiss/",
			"https://jobs.dou.ua/companies/axdraft/",
			"https://jobs.dou.ua/companies/favbet/",
			"https://jobs.dou.ua/companies/atdi-inzhiniring/",
			"https://jobs.dou.ua/companies/n-ix/",
			"https://jobs.dou.ua/companies/upstars/",
			"https://jobs.dou.ua/companies/teamdev/",
			"https://jobs.dou.ua/companies/adwise-agency/",
			"https://jobs.dou.ua/companies/kraken-leads/",
			"https://jobs.dou.ua/companies/issp/",
			"https://jobs.dou.ua/companies/trionika/",
			"https://jobs.dou.ua/companies/fluvius/",
			"https://jobs.dou.ua/companies/mgid/",
			"https://jobs.dou.ua/companies/netcracker/",
			"https://jobs.dou.ua/companies/powercodelab/",
			"https://jobs.dou.ua/companies/wargaming/",
			"https://jobs.dou.ua/companies/snap/",
			"https://jobs.dou.ua/companies/-arhivizer-arhitekturnaya-vizualizatsiya-/",
			"https://jobs.dou.ua/companies/agri-chain/",
			"https://jobs.dou.ua/companies/vymex/",
			"https://jobs.dou.ua/companies/evo/",
			"https://jobs.dou.ua/companies/youscan/",
			"https://jobs.dou.ua/companies/aurocraft/",
			"https://jobs.dou.ua/companies/traderevolution/",
			"https://jobs.dou.ua/companies/firelink/",
			"https://jobs.dou.ua/companies/trinetix/",
			"https://jobs.dou.ua/companies/frag-lab/",
			"https://jobs.dou.ua/companies/zagrava-games-by-playrix/",
			"https://jobs.dou.ua/companies/wmg-international/",
			"https://jobs.dou.ua/companies/trendformer/",
			"https://jobs.dou.ua/companies/veramed/",
			"https://jobs.dou.ua/companies/parimatch-tech/",
			"https://jobs.dou.ua/companies/gismart_com/",
			"https://jobs.dou.ua/companies/lemberg-solutions/",
			"https://jobs.dou.ua/companies/simcord/",
			"https://jobs.dou.ua/companies/computools/",
			"https://jobs.dou.ua/companies/smart-business/",
			"https://jobs.dou.ua/companies/retail-innovation-raiffeisen-bank-international/",
			"https://jobs.dou.ua/companies/infomir/",
			"https://jobs.dou.ua/companies/triangu/",
			"https://jobs.dou.ua/companies/tov-renome-smart/",
			"https://jobs.dou.ua/companies/ria/",
			"https://jobs.dou.ua/companies/cs-ltd/",
			"https://jobs.dou.ua/companies/28software/",
			"https://jobs.dou.ua/companies/readdle-inc/",
			"https://jobs.dou.ua/companies/cybridge-technology/",
			"https://jobs.dou.ua/companies/unilime-group/",
			"https://jobs.dou.ua/companies/admitad/",
			"https://jobs.dou.ua/companies/plarium/",
			"https://jobs.dou.ua/companies/de-novo/",
			"https://jobs.dou.ua/companies/itransition/",
			"https://jobs.dou.ua/companies/health-joy/",
			"https://jobs.dou.ua/companies/diligences-inc/",
			"https://jobs.dou.ua/companies/cloud-works/",
			"https://jobs.dou.ua/companies/geniusee/",
			"https://jobs.dou.ua/companies/epom/",
			"https://jobs.dou.ua/companies/buki/",
			"https://jobs.dou.ua/companies/neonomics-ukraine/",
			"https://jobs.dou.ua/companies/ajax-systems/",
			"https://jobs.dou.ua/companies/nextiva/",
			"https://jobs.dou.ua/companies/softjourn/",
			"https://jobs.dou.ua/companies/andersen/",
			"https://jobs.dou.ua/companies/obschestvo-s-ogranichennoj-otvetstvennostyu-ukr-pej/",
			"https://jobs.dou.ua/companies/bms-consulting/",
			"https://jobs.dou.ua/companies/mdfin/",
			"https://jobs.dou.ua/companies/jooble/",
			"https://jobs.dou.ua/companies/atb/",
			"https://jobs.dou.ua/companies/elinext/",
			"https://jobs.dou.ua/companies/sherif/",
			"https://jobs.dou.ua/companies/banza/",
			"https://jobs.dou.ua/companies/evo/",
			"https://jobs.dou.ua/companies/kernel/",
			"https://jobs.dou.ua/companies/flexreality/",
			"https://jobs.dou.ua/companies/mwdn/",
			"https://jobs.dou.ua/companies/lyft/",
			"https://jobs.dou.ua/companies/softpositive/",
			"https://jobs.dou.ua/companies/transoftgroup/",
			"https://jobs.dou.ua/companies/bigcommerce/",
			"https://jobs.dou.ua/companies/pricewaterhousecoopers/",
			"https://jobs.dou.ua/companies/ukrainskij-protsessingovyij-tsentr-upc/",
			"https://jobs.dou.ua/companies/echostar-ukraine-llc/",
			"https://jobs.dou.ua/companies/711media/",
			"https://jobs.dou.ua/companies/suntech-innovation/",
			"https://jobs.dou.ua/companies/tasoft/",
			"https://jobs.dou.ua/companies/datarobot/",
			"https://jobs.dou.ua/companies/pixagon-games/",
			"https://jobs.dou.ua/companies/fozzy/",
			"https://jobs.dou.ua/companies/ardas-group/",
			"https://jobs.dou.ua/companies/forbytes/",
			"https://jobs.dou.ua/companies/enkonix/",
			"https://jobs.dou.ua/companies/litnet/",
			"https://jobs.dou.ua/companies/solvd/",
			"https://jobs.dou.ua/companies/tov-helsi-yua/",
			"https://jobs.dou.ua/companies/rozdoum/",
			"https://jobs.dou.ua/companies/intsurfing-llc/",
			"https://jobs.dou.ua/companies/cybridge-technology/",
			"https://jobs.dou.ua/companies/simcorp/",
			"https://jobs.dou.ua/companies/mgny-consulting/",
			"https://jobs.dou.ua/companies/mobidev/",
			"https://jobs.dou.ua/companies/softconstruct-ukraine/",
			"https://jobs.dou.ua/companies/kovalska/",
			"https://jobs.dou.ua/companies/software-development-hub/",
			"https://jobs.dou.ua/companies/crowdin/",
			"https://jobs.dou.ua/companies/ciklum/",
			"https://jobs.dou.ua/companies/byteant/",
			"https://jobs.dou.ua/companies/doc-ua/",
			"https://jobs.dou.ua/companies/g5-entertainment-ab/",
			"https://jobs.dou.ua/companies/payris/",
			"https://jobs.dou.ua/companies/starladder/",
			"https://jobs.dou.ua/companies/solus-agency/",
			"https://jobs.dou.ua/companies/brights/",
			"https://jobs.dou.ua/companies/kitsoft/",
			"https://jobs.dou.ua/companies/aimap/",
			"https://jobs.dou.ua/companies/sap-ukraine/",
			"https://jobs.dou.ua/companies/goit/",
			"https://jobs.dou.ua/companies/lun-ua/",
			"https://jobs.dou.ua/companies/stakelogic/",
			"https://jobs.dou.ua/companies/mita-teknik/",
			"https://jobs.dou.ua/companies/quantum-international/",
			"https://jobs.dou.ua/companies/lampa/",
			"https://jobs.dou.ua/companies/infopulse/",
			"https://jobs.dou.ua/companies/ulysses-graphics/",
			"https://jobs.dou.ua/companies/edvantis/",
			"https://jobs.dou.ua/companies/amedia/",
			"https://jobs.dou.ua/companies/laba/",
			"https://jobs.dou.ua/companies/star/",
			"https://jobs.dou.ua/companies/admitad/",
			"https://jobs.dou.ua/companies/gruppa-kompanij-faktor/",
			"https://jobs.dou.ua/companies/dtek-llc/",
			"https://jobs.dou.ua/companies/tov-badm/",
			"https://jobs.dou.ua/companies/teamsoft/",
			"https://jobs.dou.ua/companies/ey-ukraine/",
			"https://jobs.dou.ua/companies/a-level-ukraine/",
			"https://jobs.dou.ua/companies/bdo-llc/",
			"https://jobs.dou.ua/companies/ubisoft/",
			"https://jobs.dou.ua/companies/globallogic/",
			"https://jobs.dou.ua/companies/visa/",
			"https://jobs.dou.ua/companies/unillence/",
			"https://jobs.dou.ua/companies/airslate/",
			"https://jobs.dou.ua/companies/datamola/",
			"https://jobs.dou.ua/companies/billtech/",
			"https://jobs.dou.ua/companies/best-run-consulting/",
			"https://jobs.dou.ua/companies/singlestore/",
			"https://jobs.dou.ua/companies/de7-partner-llc/",
			"https://jobs.dou.ua/companies/quipu-gmbh/",
			"https://jobs.dou.ua/companies/kyivstar/",
			"https://jobs.dou.ua/companies/qoders-koders-ukraina/",
			"https://jobs.dou.ua/companies/winstars-technology/",
			"https://jobs.dou.ua/companies/innovecs/",
			"https://jobs.dou.ua/companies/pettersonapps/",
			"https://jobs.dou.ua/companies/archer-software/",
			"https://jobs.dou.ua/companies/cadeau-concepten/",
			"https://jobs.dou.ua/companies/voloshyn-company/",
			"https://jobs.dou.ua/companies/edsson/",
			"https://jobs.dou.ua/companies/innovation-group/",
			"https://jobs.dou.ua/companies/autodoc/",
			"https://jobs.dou.ua/companies/beetroot-academy/",
			"https://jobs.dou.ua/companies/amdaris-group/",
			"https://jobs.dou.ua/companies/luxena/",
			"https://jobs.dou.ua/companies/preste/",
			"https://jobs.dou.ua/companies/briteam/",
			"https://jobs.dou.ua/companies/e-tender/",
			"https://jobs.dou.ua/companies/taxer/",
			"https://jobs.dou.ua/companies/health-24/",
			"https://jobs.dou.ua/companies/techmagic/",
			"https://jobs.dou.ua/companies/bpc-ukraine/",
			"https://jobs.dou.ua/companies/deloitte/",
			"https://jobs.dou.ua/companies/serena-software/",
			"https://jobs.dou.ua/companies/newsoft/",
			"https://jobs.dou.ua/companies/nordcurrent-dnipro/",
			"https://jobs.dou.ua/companies/its-group/",
			"https://jobs.dou.ua/companies/squad/",
			"https://jobs.dou.ua/companies/unity-bars/",
			"https://jobs.dou.ua/companies/liki24/",
			"https://jobs.dou.ua/companies/humancreed/",
			"https://jobs.dou.ua/companies/itip/",
			"https://jobs.dou.ua/companies/lizard-soft/",
			"https://jobs.dou.ua/companies/gigacenter/",
			"https://jobs.dou.ua/companies/viseven/",
			"https://jobs.dou.ua/companies/octal-security/",
			"https://jobs.dou.ua/companies/almexoft/",
			"https://jobs.dou.ua/companies/sensoramalab/",
			"https://jobs.dou.ua/companies/mate-academy/",
			"https://jobs.dou.ua/companies/easypay/",
			"https://jobs.dou.ua/companies/lalafo/",
			"https://jobs.dou.ua/companies/elementica/",
			"https://jobs.dou.ua/companies/am-bits-llc/",
			"https://jobs.dou.ua/companies/evolution/",
			"https://jobs.dou.ua/companies/blackbird-lab/",
			"https://jobs.dou.ua/companies/gerc/",
			"https://jobs.dou.ua/companies/automation-spirit-alliance/",
			"https://jobs.dou.ua/companies/hitica/",
			"https://jobs.dou.ua/companies/realview-ua/",
			"https://jobs.dou.ua/companies/ailand-systems-ltd/",
			"https://jobs.dou.ua/companies/nokia/",
			"https://jobs.dou.ua/companies/weplay-esports/",
			"https://jobs.dou.ua/companies/gsc-game-world/",
			"https://jobs.dou.ua/companies/enim-soft/",
			"https://jobs.dou.ua/companies/zelh/",
			"https://jobs.dou.ua/companies/smartcraft/",
			"https://jobs.dou.ua/companies/ooo-ukrnet/",
			"https://jobs.dou.ua/companies/cleverfleet/",
			"https://jobs.dou.ua/companies/ubos-tech/",
			"https://jobs.dou.ua/companies/aspichi/",
			"https://jobs.dou.ua/companies/respeecher/",
			"https://jobs.dou.ua/companies/webcase/",
			"https://jobs.dou.ua/companies/stfalcon-studio/",
			"https://jobs.dou.ua/companies/tabletki-ua/",
			"https://jobs.dou.ua/companies/watched-tech/",
			"https://jobs.dou.ua/companies/it-smartflex/",
			"https://jobs.dou.ua/companies/igis/",
			"https://jobs.dou.ua/companies/thingsboard-inc/",
			"https://jobs.dou.ua/companies/relokia/",
			"https://jobs.dou.ua/companies/eigenmethod/",
			"https://jobs.dou.ua/companies/ftech-ua/",
			"https://jobs.dou.ua/companies/novaposhta/",
			"https://jobs.dou.ua/companies/dlp-group/",
			"https://jobs.dou.ua/companies/moohii-com/",
			"https://jobs.dou.ua/companies/medirent/",
			"https://jobs.dou.ua/companies/spaceberry/",
			"https://jobs.dou.ua/companies/bit-sistems/",
			"https://jobs.dou.ua/companies/memcrab-studio/",
			"https://jobs.dou.ua/companies/elomia/",
			"https://jobs.dou.ua/companies/binotel/",
			"https://jobs.dou.ua/companies/wise-it/",
			"https://jobs.dou.ua/companies/cdm/",
			"https://jobs.dou.ua/companies/delasport/",
			"https://jobs.dou.ua/companies/binariks/",
			"https://jobs.dou.ua/companies/business-evolutoin/",
			"https://jobs.dou.ua/companies/bjet/",
			"https://jobs.dou.ua/companies/ooo-litiko/",
			"https://jobs.dou.ua/companies/prostopay/",
			"https://jobs.dou.ua/companies/ua-dynamics/",
			"https://jobs.dou.ua/companies/cosmonova-broadcast/",
			"https://jobs.dou.ua/companies/gis-art/",
			"https://jobs.dou.ua/companies/softwareone-ukraine/",
			"https://jobs.dou.ua/companies/eddy-lms/",
			"https://jobs.dou.ua/companies/keenethics/",
			"https://jobs.dou.ua/companies/metinvest-digital/",
			"https://jobs.dou.ua/companies/konica-minolta-ukraine/",
			"https://jobs.dou.ua/companies/oko-kamera/",
			"https://jobs.dou.ua/companies/dentsu-aegis-network-ukraine/",
			"https://jobs.dou.ua/companies/roboneers/",
			"https://jobs.dou.ua/companies/tentens-tech/",
			"https://jobs.dou.ua/companies/techery/",
			"https://jobs.dou.ua/companies/dobrobut/",
			"https://jobs.dou.ua/companies/igama-tech/",
			"https://jobs.dou.ua/companies/frs-herdwatch/",
			"https://jobs.dou.ua/companies/mustpay/",
			"https://jobs.dou.ua/companies/procard-processing-centre/",
			"https://jobs.dou.ua/companies/dentsu-aegis-network-ukraine/",
			"https://jobs.dou.ua/companies/velitech/",
			"https://jobs.dou.ua/companies/medeanalytics-inc/",
			"https://jobs.dou.ua/companies/tvbet/",
			"https://jobs.dou.ua/companies/design-bureau-kit/",
			"https://jobs.dou.ua/companies/aerobavovna/",
			"https://jobs.dou.ua/companies/brights/",
			"https://jobs.dou.ua/companies/ucloud/",
			"https://jobs.dou.ua/companies/eleks/",
			"https://jobs.dou.ua/companies/ubrainians/",
			"https://jobs.dou.ua/companies/one-service-ua/",
			"https://jobs.dou.ua/companies/lun/",
			"https://jobs.dou.ua/companies/vseosvita/",
			"https://jobs.dou.ua/companies/adindex/",
			"https://jobs.dou.ua/companies/n-ix/",
			"https://jobs.dou.ua/companies/lineup/",
			"https://jobs.dou.ua/companies/quod-financial-ukraine/",
			"https://jobs.dou.ua/companies/rolique/",
			"https://jobs.dou.ua/companies/nitralabs/",
			"https://jobs.dou.ua/companies/ria/",
			"https://jobs.dou.ua/companies/conscensia/",
			"https://jobs.dou.ua/companies/logics7-inc/",
			"https://jobs.dou.ua/companies/soft-generation/",
			"https://jobs.dou.ua/companies/youcontrol/",
			"https://jobs.dou.ua/companies/forpost/",
			"https://jobs.dou.ua/companies/uklon/",
			"https://jobs.dou.ua/companies/haiqu/",
			"https://jobs.dou.ua/companies/plvision/",
			"https://jobs.dou.ua/companies/artellence/",
			"https://jobs.dou.ua/companies/bubibo/",
			"https://jobs.dou.ua/companies/msg-plaut-ukraine/",
			"https://jobs.dou.ua/companies/vector-software-1/",
			"https://jobs.dou.ua/companies/systemgroup/",
			"https://jobs.dou.ua/companies/cases/",
			"https://jobs.dou.ua/companies/grid-dynamics/",
			"https://jobs.dou.ua/companies/youcontrol/",
			"https://jobs.dou.ua/companies/devox-software/",
			"https://jobs.dou.ua/companies/join-to-it/",
			"https://jobs.dou.ua/companies/roboneers/",
			"https://jobs.dou.ua/companies/reenbit/",
			"https://jobs.dou.ua/companies/dinamicka-development/",
			"https://jobs.dou.ua/companies/codeit/",
			"https://jobs.dou.ua/companies/esper-bionics/",
			"https://jobs.dou.ua/companies/controlpay/",
			"https://jobs.dou.ua/companies/melexis/",
			"https://jobs.dou.ua/companies/leobit/",
			"https://jobs.dou.ua/companies/uniqit-group-llc/",
			"https://jobs.dou.ua/companies/tango-me/",
			"https://jobs.dou.ua/companies/tabletki-ua/",
			"https://jobs.dou.ua/companies/devlight/",
			"https://jobs.dou.ua/companies/uniweb/",
			"https://jobs.dou.ua/companies/infomir/",
			"https://jobs.dou.ua/companies/provectus/",
			"https://jobs.dou.ua/companies/retail-group/",
			"https://jobs.dou.ua/companies/traffic-devils/",
			"https://jobs.dou.ua/companies/traffic-devils/",
			"https://jobs.dou.ua/companies/sombra/",
			"https://jobs.dou.ua/companies/playtika-ua/",
			"https://jobs.dou.ua/companies/merezha-aptek-bazhayemo-zdorovya/",
			"https://jobs.dou.ua/companies/quartsoft-1/",
			"https://jobs.dou.ua/companies/softserve/",
			"https://jobs.dou.ua/companies/pettersonapps/",
			"https://jobs.dou.ua/companies/netwave/",
			"https://jobs.dou.ua/companies/itera-consulting-group-ukraine/",
			"https://jobs.dou.ua/companies/zbs-tech/",
			"https://jobs.dou.ua/companies/softcom/",
			"https://jobs.dou.ua/companies/introduct/",
			"https://jobs.dou.ua/companies/introduct/",
			"https://jobs.dou.ua/companies/volmi-games/",
			"https://jobs.dou.ua/companies/lomspace/",
			"https://jobs.dou.ua/companies/muk/",
			"https://jobs.dou.ua/companies/xoresearch/",
		},
	},

	{
		Prefix: "https://djinni.co/jobs/",
		Prefixes: []string{
			"https://djinni.co/jobs/company-allright-cc765/",
			"https://djinni.co/jobs/company-englishdom-209e8/",
			"https://djinni.co/jobs/company-powercode-9f88a/",
			"https://djinni.co/jobs/company-genesis-bbc83/",                   // genesis
			"https://djinni.co/jobs/company-gen-tech-f1f4f/",                  // genesis
			"https://djinni.co/jobs/company-gen-tech2-427ee/",                 // genesis
			"https://djinni.co/jobs/company-genesis-tech-b88e2/",              // genesis
			"https://djinni.co/jobs/company-amomedia-4c317/",                  // genesis
			"https://djinni.co/jobs/company-headway-app-81bee/",               // genesis
			"https://djinni.co/jobs/company-solid-fintech-company-87d5d/",     // genesis
			"https://djinni.co/jobs/company-solid-05f2a/",                     // genesis
			"https://djinni.co/jobs/company-nebula-project-by-genesis-ed113/", // genesis
			"https://djinni.co/jobs/company-socialtech-6b80f/",                // genesis
			"https://djinni.co/jobs/company-genesis-holywater--6c47c/",        // genesis
			"https://djinni.co/jobs/company-genesis-apps-fintech-58625/",      // genesis
			"https://djinni.co/jobs/company-dc-de2ed/",                        // genesis
			"https://djinni.co/jobs/company-betterme-4372c/",                  // genesis
			"https://djinni.co/jobs/company-redtrack-io-3b68b/",               // genesis
			"https://djinni.co/jobs/company-gen-tech-growth-team--a4841/",     // genesis
			"https://djinni.co/jobs/?keywords=Genesis",                        // genesis
			"https://djinni.co/jobs/company-jooble-e95dd/",
			"https://djinni.co/jobs/company-netpeak-group-20216/",
			"https://djinni.co/jobs/company-saldo-apps-e0cc6/", // netpeak
			"https://djinni.co/jobs/company-leverx-com-47815/",
			"https://djinni.co/jobs/company-refaceai-7538f/",
			"https://djinni.co/jobs/company-s-pro-4cd02/",
			"https://djinni.co/jobs/company-banza-f738b/",
			"https://djinni.co/jobs/company-datagroup-a26a2/",
			"https://djinni.co/jobs/company-deusrobots-com-46e96/",
			"https://djinni.co/jobs/company-parimatch-tech-b6a34/",
			"https://djinni.co/jobs/company-parimatch-international-1b06c/",
			"https://djinni.co/jobs/company-pokermatch-com-43742/",
			"https://djinni.co/jobs/company-epam-systems-bb0df/",
			"https://djinni.co/jobs/company-ajax-systems-8b02d/",
			"https://djinni.co/jobs/company-kyivstar-c5f1a/",
			"https://djinni.co/jobs/company-nika-tech-family-d3d98/", // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://djinni.co/jobs/company-sigma-software-c03a7/",
			"https://djinni.co/jobs/company-ideasoft-dbe69/",
			"https://djinni.co/jobs/company-ideasoft-io-40c97/",
			"https://djinni.co/jobs/company-intetics-5221d/",
			"https://djinni.co/jobs/company-intetics-minsk-cb2dc/",
			"https://djinni.co/jobs/company-playson-05bc8/",
			"https://djinni.co/jobs/company-smartiway-6d23c/",
			"https://djinni.co/jobs/company-ciklum-international-80662/",
			"https://djinni.co/jobs/company-nix-solutions-fe08e/",
			"https://djinni.co/jobs/company-softserve-6bee7/",
			"https://djinni.co/jobs/company-softserve-dnipro-58f42/",
			"https://djinni.co/jobs/company-softserve-lviv-0de17/",
			"https://djinni.co/jobs/company-softserve-kharkiv-9b88e/",
			"https://djinni.co/jobs/company-raiffeisen-bank-aval-7b988/",
			"https://djinni.co/jobs/company-raiffeisen-bank-international-4cb67/",
			"https://djinni.co/jobs/company-eleks-6227d/",
			"https://djinni.co/jobs/company-petcube-com-a1c55/",
			"https://djinni.co/jobs/company-softteco-c9269/",
			"https://djinni.co/jobs/company-luxoft-ec4fe/",
			"https://djinni.co/jobs/company-adtelligent-751ce/",
			"https://djinni.co/jobs/company-intersog-fa680/",
			"https://djinni.co/jobs/company-gismart-e0be5/",
			"https://djinni.co/jobs/company-softpositive-37a51/",
			"https://djinni.co/jobs/company-vodafone-ukraine-85a78/",
			"https://djinni.co/jobs/company-plvision-a0d4d/",
			"https://djinni.co/jobs/company-treeum-6d9c0/",
			"https://djinni.co/jobs/company-globallogic-43eee/",
			"https://djinni.co/jobs/company-graintrack-cd726/",
			"https://djinni.co/jobs/company-intellias-c99a7/",        // question
			"https://djinni.co/jobs/company-astound-commerce-a1b13/", // https://dou.ua/lenta/articles/what-it-companies-think-about-bill-5376/
			"https://djinni.co/jobs/company-wargaming-325df/",
			"https://djinni.co/jobs/company-glovo-a9cf1/",                   // Glovo
			"https://djinni.co/jobs/company-don-t-panic-recruitment-4f656/", // Glovo proxy https://web.archive.org/web/20211111175135/https://djinni.co/jobs/289484-frontend-engineer-v-glovo/
			"https://djinni.co/jobs/company-staffingpartner-162e0/",
			"https://djinni.co/jobs/company-smart-solutions-33094/", // https://jobs.dou.ua/companies/smart-solutions/
			"https://djinni.co/jobs/company-fintech-band-3a8e1/",
			"https://djinni.co/jobs/company-fintech-farm-93fc9/",
			"https://djinni.co/jobs/company-itransition-6d0d7/", // https://dou.ua/forums/topic/35889/

			"https://djinni.co/jobs/company-sendpulse-d0078/",     // https://dou.ua/lenta/articles/companies-about-diia-city/
			"https://djinni.co/jobs/company-smart-project-9debc/", // https://dou.ua/lenta/articles/companies-about-diia-city/
			"https://djinni.co/jobs/company-zagrava-70810/",       // https://dou.ua/lenta/articles/companies-about-diia-city/

			"https://djinni.co/jobs/company-heyman-ai-08229/",  // revolut https://dou.ua/lenta/news/diia-city-has-been-officially-launched/
			"https://djinni.co/jobs/company-macpaw-76eae/",     // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/
			"https://djinni.co/jobs/company-govitall-f73ad/",   // calmerry https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://djinni.co/jobs/company-innovecs-2a027/",   // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://djinni.co/jobs/company-ilogos-c864c/",     // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://djinni.co/jobs/company-aiesec-net-e5eff/", // natus-vincere https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469

			"https://djinni.co/jobs/company-megogo-b5ed6/", // https://dev.ua/news/megogo-vstupaeie-v-diia-sity

			"https://djinni.co/jobs/company-roosh-recruitment-eac0d/", // https://dou.ua/lenta/news/44-companies-applied-to-join-diia-city/
			"https://djinni.co/jobs/company-trios-systems-56c37/",     // https://dou.ua/lenta/interviews/nuzhnyi-about-trios/
			"https://djinni.co/jobs/company-softblues-df7bb/",         // https://dou.ua/lenta/interviews/nuzhnyi-about-trios/

			"https://djinni.co/jobs/company-samsung-r-d-institute-ukraine-2c782/",       // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://djinni.co/jobs/company-samsung-electronics-ukraine-company-a2658/", // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://djinni.co/jobs/company-ria-3b2bd/",                                 // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://djinni.co/jobs/company-rozetka-54c48/",                             // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://djinni.co/jobs/company-plarium-kharkiv-c465f/",                     // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://djinni.co/jobs/company-plarium-kyiv-6e4d7/",                        // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://djinni.co/jobs/company-infoplius-b9038/",                           // https://dou.ua/lenta/news/55-first-residents-of-diia-city/

			"https://djinni.co/jobs/company-axels-tech-3885b/",                // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-playrix-63a38/",                   // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-platrix-d6928/",                   // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-span-c8e1c/",                      // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-product-madness-e9e93/",           // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-ipland-a1334/",                    // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-altexsoft-1a332/",                 // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-voki-games-6e8b3/",                // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-friendly-fox-and-4friends-16126/", // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-cqg-si-k-iu-dzhi-14fb8/",          // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-boolat-play-playrix--90ef7/",      // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company--dbo-soft--43e15/",                // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-druk-ua-050ed/",                   // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-ppl33-35-7a0c2/",                  // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-e-consulting-0b597/",              // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-pm-partners-80042/",               // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-codiv-1dfa9/",                     // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-ligazakon-ua-36e6e/",              // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-home-games-66cdf/	",               // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-jet-bi-edd8e/",                    // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-taurus-quadra-8f1e1/",             // https://dou.ua/lenta/articles/diia-city-registry/
			"https://djinni.co/jobs/company-omo-systems-aae77/",               // https://dou.ua/lenta/articles/diia-city-registry/

			"https://djinni.co/jobs/company-trinetix-58b20/", // from email letter inside 2022-02-16
			"https://djinni.co/jobs/?company=orderry-4c6fb",  // https://dou.ua/lenta/articles/diia-city-registry/#2525061
		},
	},
	{
		Prefix: "https://www.linkedin.com/company/",
		Prefixes: []string{
			"https://www.linkedin.com/company/allright-com/",
			"https://www.linkedin.com/company/englishdom/",
			"https://www.linkedin.com/company/powercodecouk/",
			"https://www.linkedin.com/company/powercodeacademy/",

			"https://www.linkedin.com/company/genesis-technology-partners/",
			"https://www.linkedin.com/company/genesis-investments-vc/",  // Genesis Investments
			"https://www.linkedin.com/company/gthw-app-limited/",        // Genesis Headway
			"https://www.linkedin.com/company/amomedia-company/",        // Genesis AmoMedia
			"https://www.linkedin.com/company/obrio-genesis/",           // Genesis Obrio
			"https://www.linkedin.com/company/solidgate-technologies/",  // Genesis Solid
			"https://www.linkedin.com/company/solidgate/",               // Genesis Solid
			"https://www.linkedin.com/company/socialtechnologies/",      // Genesis
			"https://www.linkedin.com/company/boosters-apps/",           // Genesis
			"https://www.linkedin.com/company/betterme-apps/",           // Genesis
			"https://www.linkedin.com/company/amomama-media-publisher/", // Genesis
			"https://www.linkedin.com/company/adbraze/",                 // Genesis
			"https://www.linkedin.com/company/flyer-one-vc/",            // Genesis
			"https://www.linkedin.com/company/redtrackio/",              // Genesis
			"https://www.linkedin.com/company/sendios/",                 // Genesis
			"https://www.linkedin.com/company/sendiosu-ukraine/",        // Genesis
			"https://www.linkedin.com/company/keiki-tech/",              // Genesis
			"https://www.linkedin.com/company/appflameofficial/",        // Genesis
			"https://www.linkedin.com/company/quarks-tech/",             // Genesis
			"https://www.linkedin.com/company/jooble/",
			"https://www.linkedin.com/company/netpeak-group/",       // netpeak
			"https://www.linkedin.com/company/netpeak/",             // netpeak
			"https://www.linkedin.com/company/ringostat/",           // netpeak
			"https://www.linkedin.com/company/ringostat-belarus/",   // netpeak
			"https://www.linkedin.com/company/serpstat/",            // netpeak
			"https://www.linkedin.com/company/saldo-apps/",          // netpeak
			"https://www.linkedin.com/company/academyocean/",        // netpeak
			"https://www.linkedin.com/company/tonti-laguna/",        // netpeak
			"https://www.linkedin.com/company/tonti-laguna-mobile/", // netpeak
			"https://www.linkedin.com/company/inweb-ua/",            // netpeak
			"https://www.linkedin.com/company/leverx/",
			"https://www.linkedin.com/company/leverx-group/",
			"https://www.linkedin.com/company/refaceapp/",
			"https://www.linkedin.com/company/s-pro/",
			"https://www.linkedin.com/company/banzait/",
			"https://www.linkedin.com/company/datagroup1/",
			"https://www.linkedin.com/company/deus-robots/",

			"https://www.linkedin.com/company/parimatch-tech/",
			"https://www.linkedin.com/company/parimatch-tech1/",
			"https://www.linkedin.com/company/parimatch-hq-cy/",
			"https://www.linkedin.com/company/parimatch-tech-career/",
			"https://www.linkedin.com/company/parimatch-global/",
			"https://www.linkedin.com/company/parimatch-international/",
			"https://www.linkedin.com/company/parimatch-cy/",
			"https://www.linkedin.com/company/parimatch-ukraine/",
			"https://www.linkedin.com/company/parimatch-belarus/",
			"https://www.linkedin.com/company/parimatch-africa/",
			"https://www.linkedin.com/company/parimatch-kazakhstan/",
			"https://www.linkedin.com/company/parimatch-russia/",
			"https://www.linkedin.com/company/gr8-tech/",
			"https://www.linkedin.com/company/pmint/",
			"https://www.linkedin.com/company/pokermatch/",
			"https://www.linkedin.com/company/pokermatch-ukraine/",

			"https://www.linkedin.com/company/epam-systems/",
			"https://www.linkedin.com/company/ajax-systems/",
			"https://www.linkedin.com/company/kyivstar/",
			"https://www.linkedin.com/company/kyivstar-business-hub/",
			"https://www.linkedin.com/company/nika-tech-family/",
			"https://www.linkedin.com/company/beeventures/",
			"https://www.linkedin.com/company/sigma-software-group/",
			"https://www.linkedin.com/company/sigma_group/",
			"https://www.linkedin.com/company/sigma-industry/",
			"https://www.linkedin.com/company/sigma-civil/",
			"https://www.linkedin.com/company/sigma-it-group/",
			"https://www.linkedin.com/company/sigma-technology-ab/",
			"https://www.linkedin.com/company/ideasoft.io/",
			"https://www.linkedin.com/company/intetics/",
			"https://www.linkedin.com/company/intetics-team/",
			"https://www.linkedin.com/company/playson/",
			"https://www.linkedin.com/company/zorachka/",
			"https://www.linkedin.com/company/digitalfuture/",
			"https://www.linkedin.com/company/smartiway/",
			"https://www.linkedin.com/company/smartiway-ukraine/",
			"https://www.linkedin.com/company/ciklum/",
			"https://www.linkedin.com/company/nix-solutions-ltd/",
			"https://www.linkedin.com/company/nix-community/",
			"https://www.linkedin.com/company/nixs/",
			"https://www.linkedin.com/company/softserve/",
			"https://www.linkedin.com/company/raiffeisen-ua/",
			"https://www.linkedin.com/company/eleks/",
			"https://www.linkedin.com/company/petcube/",
			"https://www.linkedin.com/company/softteco/",
			"https://www.linkedin.com/company/luxoft/",
			"https://www.linkedin.com/company/adtelligent/",
			"https://www.linkedin.com/company/intersog/",
			"https://www.linkedin.com/company/gismart/",
			"https://www.linkedin.com/company/softpositive/",
			"https://www.linkedin.com/company/vodafone-ukraine/",
			"https://www.linkedin.com/company/vodafone-retail-ukraine/",
			"https://www.linkedin.com/company/plvision/",
			"https://www.linkedin.com/company/treeum/",
			"https://www.linkedin.com/company/globallogic/",
			"https://www.linkedin.com/company/globallogicukraine/",
			"https://www.linkedin.com/company/graintrack/",
			"https://www.linkedin.com/company/intellias/",        // question
			"https://www.linkedin.com/company/astound-commerce/", // https://dou.ua/lenta/articles/what-it-companies-think-about-bill-5376/
			"https://www.linkedin.com/company/wargaming-net/",
			"https://www.linkedin.com/company/glovo-app/",            // Glovo
			"https://www.linkedin.com/company/dontpanicrecruitment/", // Glovo proxy https://web.archive.org/web/20211111175135/https://djinni.co/jobs/289484-frontend-engineer-v-glovo/
			"https://www.linkedin.com/company/staffingpartner/",
			"https://www.linkedin.com/company/smart-solutions-hr/", // https://jobs.dou.ua/companies/smart-solutions/
			"https://www.linkedin.com/company/fintech-band/",
			"https://www.linkedin.com/company/fintech-farm1/",
			"https://www.linkedin.com/company/itransition/", // https://dou.ua/forums/topic/35889/

			"https://www.linkedin.com/company/sendpulse/",            // https://dou.ua/lenta/articles/companies-about-diia-city/
			"https://www.linkedin.com/company/sendpulsebr/",          // https://dou.ua/lenta/articles/companies-about-diia-city/
			"https://www.linkedin.com/company/sendpulselatam/",       // https://dou.ua/lenta/articles/companies-about-diia-city/
			"https://www.linkedin.com/company/zagrava-games-studio/", // https://dou.ua/lenta/articles/companies-about-diia-city/

			"https://www.linkedin.com/company/revolut/",       // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/
			"https://www.linkedin.com/company/macpaw/",        // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/
			"https://www.linkedin.com/company/calmerry/",      // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://www.linkedin.com/company/govitall/",      // calmerry https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://www.linkedin.com/company/innovecs/",      // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://www.linkedin.com/company/ilogos/",        // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469
			"https://www.linkedin.com/company/natus-vincere/", // https://dou.ua/lenta/news/diia-city-has-been-officially-launched/#2343469

			"https://www.linkedin.com/company/megogo-net/", // https://dev.ua/news/megogo-vstupaeie-v-diia-sity

			"https://www.linkedin.com/company/weareroosh/",    // https://dou.ua/lenta/news/44-companies-applied-to-join-diia-city/
			"https://www.linkedin.com/company/trios-systems/", // https://dou.ua/lenta/interviews/nuzhnyi-about-trios/
			"https://www.linkedin.com/company/softblues/",     // https://dou.ua/lenta/interviews/nuzhnyi-about-trios/

			"https://www.linkedin.com/company/samsung-electronics-ukraine-company/", // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://www.linkedin.com/company/ria-com/",                             // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://www.linkedin.com/company/rozetka/",                             // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://www.linkedin.com/company/plarium/",                             // https://dou.ua/lenta/news/55-first-residents-of-diia-city/
			"https://www.linkedin.com/company/askod/",                               // https://dou.ua/lenta/news/55-first-residents-of-diia-city/

			"https://www.linkedin.com/company/axels-tech/",                                  // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/playrix-entertainment/",                       // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/daily-magic-productions/",                     // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/span/",                                        // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/product-madness/",                             // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/-ipland-llc/",                                 // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/vizor-games/",                                 // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/altexsoft/",                                   // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/investment-intelligence-ltd/",                 // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/voki-games/",                                  // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/fourfriends/",                                 // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/cqg/",                                         // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/boolat-play/",                                 // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/druk-ua/",                                     // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/%D0%BE%D0%BE%D0%BE-%D0%BF%D0%BF%D0%BB-33-35/", // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/e-consulting/",                                // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/pmpartners/",                                  // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/codiv-ukraine/",                               // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/reardengroup/",                                // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/liga-zakon/",                                  // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/showcase/ligazakon/",                                  // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/home-games-ukraine/",                          // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/perevaga/",                                    // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/tq_ltd/",                                      // https://dou.ua/lenta/articles/diia-city-registry/
			"https://www.linkedin.com/company/omosystems/",                                  // https://dou.ua/lenta/articles/diia-city-registry/

			"https://www.linkedin.com/company/trinetix-inc/", // from email letter inside 2022-02-16

			"https://www.linkedin.com/company/orderry/", // https://dou.ua/lenta/articles/diia-city-registry/#2525061
		},
	},
}

func Verify(content []byte) Response {
	var request VerifyRequest

	var unmarhsalErr = json.Unmarshal(content, &request)
	if unmarhsalErr != nil {
		return Response{`{}`, http.StatusBadRequest}
	}

	if len(request.URLs) == 0 {
		return Response{`{}`, http.StatusBadRequest}
	}

	for _, ss := range outsideWarCompanies {
		for _, requestURL := range request.URLs {
			if hasAnyPrefix(ss, requestURL) {
				return Response{unsafeOutsideWarCompanyMessage, http.StatusOK}
			}
		}
	}

	for _, prefixGroup := range stopDiiaCityPrefixes {
		for _, requestURL := range request.URLs {
			if strings.HasPrefix(requestURL, prefixGroup.Prefix) {
				if hasAnyPrefix(prefixGroup.Prefixes, requestURL) {
					return Response{unsafeMessage, http.StatusOK}
				}
			}
		}
	}

	return Response{safeMessage, http.StatusOK}
}

func Prefixes() []PrefixGroup {
	return stopDiiaCityPrefixes
}

func hasAnyPrefix(prefixes []string, url string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(url, prefix) {
			return true
		}
	}

	return false
}
