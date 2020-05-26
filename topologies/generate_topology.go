package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// keys can be generated with TestPrintGeneratedKey in common/sign_test.go

// 64 keypairs
var generatedKeysServer = [][]string{
	{"a77000e6e7e77df500e39144e34e31f1715c98a5ecaaf34b8863f2d907c35ab4", "26665d594658641d981cb13e74e130c12bb5a0924c217e10b2a0e20437cabf97a77000e6e7e77df500e39144e34e31f1715c98a5ecaaf34b8863f2d907c35ab4"},
	{"c8766183072ad3a0aeb17e672152bff7d5382a62f14a818f2a7dbe3c123d1fc3", "19d771ac0e1f2b76082820840dd8b7df7c7e2614f8d7346186460946b5f6173dc8766183072ad3a0aeb17e672152bff7d5382a62f14a818f2a7dbe3c123d1fc3"},
	{"564c67c3cd08863a570e632c60409d57b02f576a0582bcaf297aa4610f98aead", "eb1a60b82e72545395f8ce4334703f13597985aed204107ea5bdaf6289c38a68564c67c3cd08863a570e632c60409d57b02f576a0582bcaf297aa4610f98aead"},
	{"4a155d1c1d86854aae1008373c3e9a92dd02cbbef52e5372aee13c2efe22398d", "e76b8dee9a5eca0000da56aadaa39e8cdb564ee6f82008052caad39741e72b134a155d1c1d86854aae1008373c3e9a92dd02cbbef52e5372aee13c2efe22398d"},
	{"da67cdf6dcf621da70f2a3e902b371853b64d2da457de81eff5c9cb1da811f09", "fc2d5e43dea1ce49885ad705abfa0c4370d103b2c5db566b8f5903b0b7acb3b1da67cdf6dcf621da70f2a3e902b371853b64d2da457de81eff5c9cb1da811f09"},
	{"e1a4cca58d70ac85d567358d01878648abdcdb9a92f0fe289c116b27181b04b6", "3823a36ed2057ebb2b00da9c75677fca4be8829e4ecc4d55bc3b1196f8c5fbb9e1a4cca58d70ac85d567358d01878648abdcdb9a92f0fe289c116b27181b04b6"},
	{"30d615e6962bcfd30847a8aa6ed103aba903603b0f0b11a090b344aca881a114", "3b5495f1843e199b1c559320603eeb0e05e4b3ebf48d64aebd19bd7a5ce5d38230d615e6962bcfd30847a8aa6ed103aba903603b0f0b11a090b344aca881a114"},
	{"a23370ed4fdb9219e03dd354f3aa76400b440ba70e40f95615c2577c4d88e4c2", "b72ebb6b81268c9c103dd432477edbb580d9e936ad3978d6e5f037337583eaa8a23370ed4fdb9219e03dd354f3aa76400b440ba70e40f95615c2577c4d88e4c2"},
	{"a3f70dbb0a14afa7a33358bae945293dbbe355068efc90e9ff3d7dfe45136a27", "de4ce7be9625ac97c56aeb64c8c62be83f36f5898fd9b1bf72d56c1bcb7d18dfa3f70dbb0a14afa7a33358bae945293dbbe355068efc90e9ff3d7dfe45136a27"},
	{"618a1a777293968592084a0b5ef0ef6117358025fa513c978197ad027273b1f2", "877ea07f401c16c6c4b22b0fb8f38e4fd9a7766a57daca7c21b07ccab8812573618a1a777293968592084a0b5ef0ef6117358025fa513c978197ad027273b1f2"},
	{"b9d05f72b918fe07c11bb6ff3ce27294c7d281c163107e38a1676e195749aa39", "cd4d74cb539d0695d2afc892c6cc5189cf1404c51e927d9bb850002b4e7c826ab9d05f72b918fe07c11bb6ff3ce27294c7d281c163107e38a1676e195749aa39"},
	{"848af695ee0b7fc4a8be73ddc63e9fed6a627fd9a97e659f403ea30dc3d5742b", "bba5990380e0ee7b458d9b2ef322ed0cc248f452cbf66032c82778851b8e1b51848af695ee0b7fc4a8be73ddc63e9fed6a627fd9a97e659f403ea30dc3d5742b"},
	{"9610201840ac0d215eb085c51ee07d679421c4a02d3810ab8d60273825de8930", "13cb2e4b6e744cfba1347a65c3af135252926b6832274356298437e89e80cc179610201840ac0d215eb085c51ee07d679421c4a02d3810ab8d60273825de8930"},
	{"3166cf23ace83a5a39d4e21ebbfbe873c9a270623c1855da74fafed7decd3d82", "7946cc0a057a517667ea49029783f2d93ecdf5a4383f5792ac76d6a23a3c3e083166cf23ace83a5a39d4e21ebbfbe873c9a270623c1855da74fafed7decd3d82"},
	{"c7ba8fb18fd8e4fe1ad2eda326b5367fff6d15b8f563f867c7d5db516f0930aa", "e673216259419edbb9d63ea1694190c63666b656f01d55d09adb7219dada4aa4c7ba8fb18fd8e4fe1ad2eda326b5367fff6d15b8f563f867c7d5db516f0930aa"},
	{"bbd7b5ded1c667ce3b0ef3b9602a18c9d3ff5e24e4e27e4701015a2ef91e566e", "d768c7f5e5edeb34bdf0948c0fe208fd6ff464ff91f09f5e0f2bac1609374e9dbbd7b5ded1c667ce3b0ef3b9602a18c9d3ff5e24e4e27e4701015a2ef91e566e"},
	{"0107aed1aac615c000ab8c02c13c0d1006b063c5571ac07881054f9c31bf70ed", "d0ed71505e666cb2204de4dfacd0e47cfe5c332c489c841a27fdaa4f466cfc9d0107aed1aac615c000ab8c02c13c0d1006b063c5571ac07881054f9c31bf70ed"},
	{"6ce0382756b567e79b4c2ba47e4425c263312bd80f3c9dd716c29722e894febe", "32af4ae57673464982edfa4a06ad8b9c917729a14e344dc4799553e79b31f21e6ce0382756b567e79b4c2ba47e4425c263312bd80f3c9dd716c29722e894febe"},
	{"69b11d6b9cc6599be22f49d77a2fbf31963c77a61e94c462e279b3ea0564243f", "ba89dbe703d338ca4fddb609405e935c66ee29b8d656edabda20a95200e2670269b11d6b9cc6599be22f49d77a2fbf31963c77a61e94c462e279b3ea0564243f"},
	{"12190d3c264367fcf3b4c74cfa920b98975f8bea4ca30f46334e591706cb9f4d", "9323ffdd96a773ef2e781d71c371ec9ad81a0fe99eb637de5469c42e8481fbe812190d3c264367fcf3b4c74cfa920b98975f8bea4ca30f46334e591706cb9f4d"},
	{"d004d72923401c07651a1e3e2949ba850b2d947e3dc540b50bf3d14cbd620ab0", "af662acd7e461484c6bca9ae3b5ba43e846c4e62c8d48fc4db1086f58f921d1dd004d72923401c07651a1e3e2949ba850b2d947e3dc540b50bf3d14cbd620ab0"},
	{"d698db98663191fb65684e1a2cffc4c6c261925e33d4bcea0db62967b86fa39a", "79a2d469e59dfa9040f1772d427ce77f401cd2c044ae08891a5403d324c2e509d698db98663191fb65684e1a2cffc4c6c261925e33d4bcea0db62967b86fa39a"},
	{"373c77d650a74492d58c290e702a145c25ea73416f599e45df870cf144e39b95", "b3ca0d0eb80ff1834964f09e246957110643852136df08e7fcaff799b11cc055373c77d650a74492d58c290e702a145c25ea73416f599e45df870cf144e39b95"},
	{"cc283b87078e19aee44fa30b516bd442f037bee0b94480f23e78312512cb3337", "64ff2de27d7d05c9e1e47fd2aa87b88aac66c6e350bc239720688904c487fea9cc283b87078e19aee44fa30b516bd442f037bee0b94480f23e78312512cb3337"},
	{"30e8e32bfdef83879adf45f4f95be3463fc03cc79729adf16ba0442f859e924f", "c74c3837ff205a0451daa8ac67bababa016cfa8b6992c1e7b6e737efb2b8642230e8e32bfdef83879adf45f4f95be3463fc03cc79729adf16ba0442f859e924f"},
	{"0cc54486d4ca53c823deb96cc7ff0defd960a478f98f9f5e113ffe94927e15ee", "2d1faf7970939ccaf0459afe4fd34d131380434e3b64a0faf2474353cec770700cc54486d4ca53c823deb96cc7ff0defd960a478f98f9f5e113ffe94927e15ee"},
	{"9a6139d0287fe0b65d5c54e41afed68d17d6e6fcc5aeb122402b2b902cea3589", "5647af260405e7197e28aa3c0f726c5f754e29d9b864343ca0b6c174db1f0fc09a6139d0287fe0b65d5c54e41afed68d17d6e6fcc5aeb122402b2b902cea3589"},
	{"e5b20db307f5b5251b038ae5f0a8019867e27bafb928ec02a4498f42c3b0a7a1", "42c50aaa7b7ad1e40ed8f7101833638b5edec66c59bb7f41627c688dceba69fae5b20db307f5b5251b038ae5f0a8019867e27bafb928ec02a4498f42c3b0a7a1"},
	{"ecb5cf84cf08a80a5d7d44c1f125f19b6db8a8e6693b492e6713bf1cac5a3950", "523c0be425b9204ad16341a35603a1478e5f7a9db8a28e0d64e4da2dc29456d8ecb5cf84cf08a80a5d7d44c1f125f19b6db8a8e6693b492e6713bf1cac5a3950"},
	{"968ff33e6c63cbd241797c1edd28b54e495ec0c0b249ec89b3d669e240a43a0d", "a3b37111a2e72e648271624c9f79a402bd24907bde02e6ec258be364d14ef3dd968ff33e6c63cbd241797c1edd28b54e495ec0c0b249ec89b3d669e240a43a0d"},
	{"a00ad9edea172694306266b0aabd80f145cfc7d84755c23de5907a630ec787ad", "3e6aa2278db0220974b33a7122a4ce837a6238321e9db950c419e89a02f8c9d6a00ad9edea172694306266b0aabd80f145cfc7d84755c23de5907a630ec787ad"},
	{"cf1cf517f580e2a05fdba5e9e21b4477e8e87c4a652fcaec6ab32f2a3bbe0859", "cadfbb8e7bb5514f44224458de482a27ad6c192cacafd436e9aa4ee1f3f87080cf1cf517f580e2a05fdba5e9e21b4477e8e87c4a652fcaec6ab32f2a3bbe0859"},
	{"cfff495f863ddf35ccc64dc41608d04cf6d9dea374a466682a2a72c4592fc681", "42f76d18b1417cdc03b953edac34ca352f52c96b967343f3b28d6bf0605efb42cfff495f863ddf35ccc64dc41608d04cf6d9dea374a466682a2a72c4592fc681"},
	{"3fcc0de33db036048a91f1dee69eec0ee34e16c64554ba046055ce2e76bca3f3", "ea6fbd0597864e6c4f20f70323c01aec95b376fca9e06a5a599d7deb3ba438453fcc0de33db036048a91f1dee69eec0ee34e16c64554ba046055ce2e76bca3f3"},
	{"cbed7683ef0152fb144d41165b3740740823402cef481d22d07791f6faf0d135", "ae69e9d5095a1857dc700d96a5a4d40a555e6fbc4529fab00eaf8d14e4448f88cbed7683ef0152fb144d41165b3740740823402cef481d22d07791f6faf0d135"},
	{"619d87c2aa11f712c8b05e1d3659232c45fde1ddde247a02bd8c750ff5761a6f", "5068109926395a9be27ea61956af7c88f4594320a61a8317df2da4268f8573ad619d87c2aa11f712c8b05e1d3659232c45fde1ddde247a02bd8c750ff5761a6f"},
	{"d4f29dda0bcfc1dd213fe2a2d611b4f00f60c7dee193c0842f1a3f9c08b44569", "42336cc03ef722f4f2c1faaafc7d3b1a1e9a9ce2d02ece7e8bed28a42370d2c5d4f29dda0bcfc1dd213fe2a2d611b4f00f60c7dee193c0842f1a3f9c08b44569"},
	{"07efd1cfe0e092811ca0524d47bae4446c004aeddaf88d2bec363d56592a4b65", "f17b53093062906e0d0950c629e00eb047e765b5bb4e6ce6a72677f64eba1aa507efd1cfe0e092811ca0524d47bae4446c004aeddaf88d2bec363d56592a4b65"},
	{"55dc6c5fca03536392cbcb55f5f81aeb9fc20027e15ae0c344241655039feb8d", "deee292b6a18bcafc1930c67bc32795633b615e13a7a9bc9482746e7155b348355dc6c5fca03536392cbcb55f5f81aeb9fc20027e15ae0c344241655039feb8d"},
	{"895ffafe427ab6ebb08f286a73a5f7b2a81973307c36e0db08b47bb8442b6502", "7fc16b2448b8232d49ce57306dea5a3be8a2adbb5e42ae4815f1aff257660994895ffafe427ab6ebb08f286a73a5f7b2a81973307c36e0db08b47bb8442b6502"},
	{"9be15bee6c03339583c45034228c97178b9118c99cbb4e008b35c02c38741e2a", "d9acc77c6d3fc2f9dcd035d95d66d335759c85978154390cf1d5f7bae04858539be15bee6c03339583c45034228c97178b9118c99cbb4e008b35c02c38741e2a"},
	{"5614b5b10bf6923f24ea606acf73e46c5d4ffd3cea573d98da34d86c5114ece1", "9f747e1a4fb312acabea24301147427a924438b3d7790c35c91ba100135063025614b5b10bf6923f24ea606acf73e46c5d4ffd3cea573d98da34d86c5114ece1"},
	{"37257ef29d4d102ddb04e86ade65fc943cbb738ee0bd60b22f08ccf951c1723f", "c34ed4ccd8ad256448cd66b2d0a2c27dd54d6d466a2a2c3b7fde502ffeab5be737257ef29d4d102ddb04e86ade65fc943cbb738ee0bd60b22f08ccf951c1723f"},
	{"8a325615a9bcefe9c068f0f2a79336a1e7df05ad12d149040001d20a87679b74", "7b827beca7977fe249884f8d69f41758ba66c014931eb40b05afd66c706e56f48a325615a9bcefe9c068f0f2a79336a1e7df05ad12d149040001d20a87679b74"},
	{"8794a397da6a6e16df8648aa4f2f1b7bc922be7fcb521ffc08dbaaba8063e6ff", "14b1a6746447dfc795b9b1082912ba9031cb61f4140bc6317eeaade7a24173908794a397da6a6e16df8648aa4f2f1b7bc922be7fcb521ffc08dbaaba8063e6ff"},
	{"4bd08df8ad5f77f3c45a93605d22f417143a745b549c9b96040f84b983286dad", "4657ba9e19229b71c1fa8114226537bab808f04d4aa62d2a19b11cedbfa839c24bd08df8ad5f77f3c45a93605d22f417143a745b549c9b96040f84b983286dad"},
	{"ca3c276b9a338f64660416528c9f498d797e0bafb619d30f723f1d8a87fc6f88", "b685d70171c9f0645dcea0661805e6ef96b1e0fcf3ced36ba05668f1f48d82c9ca3c276b9a338f64660416528c9f498d797e0bafb619d30f723f1d8a87fc6f88"},
	{"f331ed896f9af2aa2c4ef32ad68d83156f61109c29c562771525827d74f7eb7a", "bd7b9a925cb23ae54da38bf479cae87b8ae98b6ce6dbca6bc670ac6abdaa9c10f331ed896f9af2aa2c4ef32ad68d83156f61109c29c562771525827d74f7eb7a"},
	{"62bd06884d76bbc6f752742f3503b72ab3be575052c27cfbac3c941fa007807e", "b34d744445f30b648aa16421547197c280039981e21fff6bc667a51bdf92769a62bd06884d76bbc6f752742f3503b72ab3be575052c27cfbac3c941fa007807e"},
	{"4d4f67d3a5fd74dd560feeeb8962a690031b1113f9bcad4c8ee40e930f004e9c", "4ea2b766237708a6f93710759917959fa67e44b3478cabe48d37511b8c2627e54d4f67d3a5fd74dd560feeeb8962a690031b1113f9bcad4c8ee40e930f004e9c"},
	{"d535b8f5c8a40478514ed046c2ee9ea5d7ddd06303d82278b0f11b81771af455", "2382610c47e72ca0464179259212e14ee2477b6e5f5715f33d131df5c3ce2a98d535b8f5c8a40478514ed046c2ee9ea5d7ddd06303d82278b0f11b81771af455"},
	{"5b97d7c9e3028976d0a55bde330abf0cbfcf3de48ef6ebaf04683dcfb644b676", "4fdaf285d528af2b2215f3046b0c8828c1518d0473dc094bb36e0a38e72fc8e75b97d7c9e3028976d0a55bde330abf0cbfcf3de48ef6ebaf04683dcfb644b676"},
	{"18658d363556a15181aad75bf06b41c3e817e6e974c12aba7e412e631412cd58", "692a628070c4ec500574f29224272de9e7b4e187e33eeed20c44a1a2bf95af3d18658d363556a15181aad75bf06b41c3e817e6e974c12aba7e412e631412cd58"},
	{"f5232af493414f26686a3c72e5562f040874c1922bd03b74fb3c254cdc2235bb", "87249016873683c0e3586ad8225f68c20260e890aa582e0a3627b247493ef8e9f5232af493414f26686a3c72e5562f040874c1922bd03b74fb3c254cdc2235bb"},
	{"28d8455cda2d184c0fedcfb0ce2a60abbd692e5fbbbc67364d0009b4da6699ec", "28a708520beb1baf957ad0e093d3d61b59d2e9ef0d48c513dbc2df971a4c156428d8455cda2d184c0fedcfb0ce2a60abbd692e5fbbbc67364d0009b4da6699ec"},
	{"0ab2154eeb702872ae6c948b44da34b7be0a18683e83685200074d8b8d3f4d7f", "291475d4c6c1939ea72de6775bb1fb9d8ec23066a0ba3481c656d3b1d135fcc40ab2154eeb702872ae6c948b44da34b7be0a18683e83685200074d8b8d3f4d7f"},
	{"75d983dbfb25a488cf3508c06272be9fd6b28eb0ef5e48765936a00fa9dcda3c", "ca9e7c50f7bcd8a2cc3ce2a533785d04a268ef3e7de993aa42b3bbb59c74963375d983dbfb25a488cf3508c06272be9fd6b28eb0ef5e48765936a00fa9dcda3c"},
	{"e1be5c462ac9896024e8863703a47d11475408f87aa6c1f81c55fe76abe6083a", "af14ea15e3def2a19031c2f1c452aa3cb43f0cb2cdd64c6da0f5a2d3601c3b37e1be5c462ac9896024e8863703a47d11475408f87aa6c1f81c55fe76abe6083a"},
	{"1122963a9d34ee9db94f34f84f7452919a6691dff0389b955130a2c03704376e", "f05e9465f0929e9374b6258d3cbd65b452b1ae5f8c6a3933b49469d1c65d404a1122963a9d34ee9db94f34f84f7452919a6691dff0389b955130a2c03704376e"},
	{"18f238142dfddb53264aec9cb9aa1327284881cf7c8935daf7ab02a3c480782a", "8ff07ce940cf3930a4943cbab29f430ca55a3539a750094945f529a9c1182eac18f238142dfddb53264aec9cb9aa1327284881cf7c8935daf7ab02a3c480782a"},
	{"c948880dbadee37e3fe2dc0f328f0753296c38418083993e3babf75ae7acde7b", "043e085a544eb65a8065397c5fa755b5078f4b0b655860c167a39d8f23d6a215c948880dbadee37e3fe2dc0f328f0753296c38418083993e3babf75ae7acde7b"},
	{"e4b4f494748a5e084d86bb72911be88809fca0db9d1249246676316e38f69407", "50b824bbaa6fd1f78fcf900ad8af1b2c4d772e44c79475f119d6bd65c686a5c4e4b4f494748a5e084d86bb72911be88809fca0db9d1249246676316e38f69407"},
	{"3430f66709bb20e6e2812191722490df6c6d8a5fe6f2b3f13a8a118af7efde39", "2b4e43df2dbbda8daae522f795ee69a86fa9950535c5ee0254e5d43ec5fda15b3430f66709bb20e6e2812191722490df6c6d8a5fe6f2b3f13a8a118af7efde39"},
	{"fcfaaca17fe89b1a18dac3d34707ee33226507839485921a26a36baefef924de", "3336d3b0b7d8aa292c0a399d094c4911b5fa54146d19e35721ea425b6a7eafb1fcfaaca17fe89b1a18dac3d34707ee33226507839485921a26a36baefef924de"},
}

// 64 keypairs
var generatedKeysClient = [][]string{
	{"5e7da91d8083e79d742f39a1e218c92fdfb8e5684b60c3d2d82bce1fcc2c8084", "135f014fc76ca27c3be0227ac8c356ab7c6cb38dda1f350bb8fda7487f8aa2dd5e7da91d8083e79d742f39a1e218c92fdfb8e5684b60c3d2d82bce1fcc2c8084"},
	{"13a619eba72c54e669d67b910d3e927f3dc1bb02764b69b4163345d55c2eb0aa", "43ccf68f7190d1fc0f6fa6b2dc35e1244290c816cd0eb4fc12d73a5c1914a86e13a619eba72c54e669d67b910d3e927f3dc1bb02764b69b4163345d55c2eb0aa"},
	{"826f9850a0685f79877b8049f3216a35b4a5b9894d4f70b50946b254ea878e22", "5f4939311ce1bc5cd889e9d35a7ea7c2db64db0d5e22eb8d8914b044ed388141826f9850a0685f79877b8049f3216a35b4a5b9894d4f70b50946b254ea878e22"},
	{"d46a71d600ad623d21ced5c5398984def57875fe4ab6d9b5d4ed6883d06e6884", "4e6eb64343c2a3c9f5241771484545fdb1bcc2749225b87d9589034c1bb81414d46a71d600ad623d21ced5c5398984def57875fe4ab6d9b5d4ed6883d06e6884"},
	{"6c8c77537d8640d7b530a4695310151bbd8b42d45f9234be29ee83b3e9419474", "ac8a1e63bb6440fd4937849cbe9296a9a92925ff19bd7e0f1f29500887a682046c8c77537d8640d7b530a4695310151bbd8b42d45f9234be29ee83b3e9419474"},
	{"cf6d3ac66bec3408012f736463590c62739798036681cd73be39071ed00a2b67", "7f8f40262ab33bfd5a973622e046fab042123b88a2114fcb97c3e6858c4b1099cf6d3ac66bec3408012f736463590c62739798036681cd73be39071ed00a2b67"},
	{"00bb932e1c935c102cf1d53452b0f10ff9ece1264ec96a6ace30bb1037ea5507", "827d7b276e4ff3efbc3d3957f156a2a958b3deae890277ab7e1a847a0415dccf00bb932e1c935c102cf1d53452b0f10ff9ece1264ec96a6ace30bb1037ea5507"},
	{"62a0061c9ec6b2a1f186dfc0f490a7eb158b072ba2d7ef049248b1337917a817", "163f4cbafb36ee4a925e22e373153e01d807a5f7de4ede92e80b93a5a4fa18c362a0061c9ec6b2a1f186dfc0f490a7eb158b072ba2d7ef049248b1337917a817"},
	{"3f4d44a67deccaa314dd203105ff9231f382a84561013cf253352059b90b692d", "e650af1698bc8299268788a9ae97935f21e970180b4f7f1a1d94e3c550e736fb3f4d44a67deccaa314dd203105ff9231f382a84561013cf253352059b90b692d"},
	{"4d886b2926a41ef39fd990c90b298cafbddd811cca2c8e678a0066c44bf31523", "faf4c5c35216b0209957b098b779cc6da9a8dbe0106556260b3035cce8b175124d886b2926a41ef39fd990c90b298cafbddd811cca2c8e678a0066c44bf31523"},
	{"bc4e0b9464a805d958d1e09e4949b451e191cf25924c52a6bed010df5d857282", "679bde789d45515b076a68517bdfbfba042f55dad50df3043b14a8384f45006fbc4e0b9464a805d958d1e09e4949b451e191cf25924c52a6bed010df5d857282"},
	{"be89288412fd3d03c0b6f512d1e5bbbc91957bf86040e0959b200933825a9e1e", "22b59214fa4fcaaadb4ffb01bd3d413e6fdfb93faa21b27f241e64c4001137efbe89288412fd3d03c0b6f512d1e5bbbc91957bf86040e0959b200933825a9e1e"},
	{"278b81feb12c387a49fadc1649e66690dd7d92e8b0419f16c1eca80b81e05e8a", "64e4ccdd76cd951aa179d32741450c925d5ee1eb70241070798badabff86dfed278b81feb12c387a49fadc1649e66690dd7d92e8b0419f16c1eca80b81e05e8a"},
	{"45b90f249fef00684a0c1ca2a395f15c8c3aedde9691b52e9ee673b547265bf5", "51d929f193bedcf1edffa72825afe4cb7757a18c81cadc7b53b9d26b2759e39445b90f249fef00684a0c1ca2a395f15c8c3aedde9691b52e9ee673b547265bf5"},
	{"26ae03b17e85cc8bee2274310c6b36ab916709ee5c93f8ae3c67b4507b48569f", "0dd90e6463ec8cab750b414f9c1d79bf7eef9cefb65304457f1abaac9f2ab64126ae03b17e85cc8bee2274310c6b36ab916709ee5c93f8ae3c67b4507b48569f"},
	{"a023dd2de6a360810033baf18ad1ea8a513bca533a67bc9b572eb3e7021dbc84", "0de3ea14676fa097d219598d062224eb374b04675e68a5dd5b778261691545fca023dd2de6a360810033baf18ad1ea8a513bca533a67bc9b572eb3e7021dbc84"},
	{"7de766dea9e1b185ba67ff3c32cbc4bcce51215a67efbdd014e2fab64e1a6393", "d9486b53b3a2b34a730fe0d09a92ab87c714c743726e9c24244bb7c75bc191587de766dea9e1b185ba67ff3c32cbc4bcce51215a67efbdd014e2fab64e1a6393"},
	{"b41dac00e557000458a8f5aa4037d5468803d65399a69ef260d4dd412c1e7837", "b8432ee6390584510bf79831c304a6d5d004818cfa7ce967b1d6e075c1b139f5b41dac00e557000458a8f5aa4037d5468803d65399a69ef260d4dd412c1e7837"},
	{"97f756e550056b384005b84e67bce22a23d0f8af76d5b8cdd54cf461f299e997", "b755fce9847cbc98761291f5c7742335bb2b1421ec813f37f1fa858babcbd11997f756e550056b384005b84e67bce22a23d0f8af76d5b8cdd54cf461f299e997"},
	{"7905bd6d35f303b5c13f3dcb62c6f9a535835948a154475d78a8d1a593f968a3", "f250231f94f52e235af624f5724fdba6498842bf9742c2375eca89d6a3bf7a7a7905bd6d35f303b5c13f3dcb62c6f9a535835948a154475d78a8d1a593f968a3"},
	{"9bd9adee1cfdc38700e7c59757e44374da1fd1053cd35488cf721ccbde6778dc", "279e2b83908a4b5523309576d7e7df6f77dab6e4832b6c2d861e6c80342bac0d9bd9adee1cfdc38700e7c59757e44374da1fd1053cd35488cf721ccbde6778dc"},
	{"eda1cda449273e5c7fe1b588d5575337a14f3baf8d9b125e5f4450639f77dde8", "318bc06bc7f85e2fac6fd17a00882d3abf79fdecd80461108241f420748a4956eda1cda449273e5c7fe1b588d5575337a14f3baf8d9b125e5f4450639f77dde8"},
	{"54dd564cf4c508ece9baa062dc71525df8ab6016047c637f16492149bac7d6b5", "ff0e63da4bc73d2055e8e1ff5e52fff0a9dec4ea27d882357301c87429f2bd7454dd564cf4c508ece9baa062dc71525df8ab6016047c637f16492149bac7d6b5"},
	{"ae182e2b3842b0e1bb1d46d479b5b06f3018f4fd547cf1b6a1e6300c021cc364", "73c6180f40c5195ca09fa20f5c9252d2e057ec9ab2fd9835440dee67f90482b1ae182e2b3842b0e1bb1d46d479b5b06f3018f4fd547cf1b6a1e6300c021cc364"},
	{"92069e8f24ca7c7ffe500d07819da64499610aad16a1cb26492777029f153b41", "78984659d7a99a478021c369b5b40b0896fab43eda12c889926340052663b1df92069e8f24ca7c7ffe500d07819da64499610aad16a1cb26492777029f153b41"},
	{"9092ca9cb3e3f38bdf5b745e6722027ddd6f247a60727f74eee6b303da445d70", "693447e3ebdee6352684d77a826cfd03be266bcbd8a7fde9a0766ef555edef829092ca9cb3e3f38bdf5b745e6722027ddd6f247a60727f74eee6b303da445d70"},
	{"68392197b8dbe8be15ced6b821a7db9c1e238870f3bc1e1daa31f91e389cac48", "9c03d73494a503907c6bbfb9e2ce886e5ad27da89540400b8df05a6e5aa0f88368392197b8dbe8be15ced6b821a7db9c1e238870f3bc1e1daa31f91e389cac48"},
	{"92274a0656060b00088cb26fde87df8b4e6365398eef3c659ee6015dd675e3eb", "52109d49e0eedc5817d18965bea37b7a78fd4c6550dfbb506a1f7264f31b075392274a0656060b00088cb26fde87df8b4e6365398eef3c659ee6015dd675e3eb"},
	{"93a032beb914cc30b922990367ade4a1c886e56f72ddc76a7e5e58f985ec888b", "11544c8e5827b45a3963634e6dd44270889626305105f22b3ab1e576f404f48993a032beb914cc30b922990367ade4a1c886e56f72ddc76a7e5e58f985ec888b"},
	{"b4f5e7b6d9dd1824391ab2d3a7593eb940428bf8c53a837236adaa9148891720", "affef99ce6236ad46525e3ac4b9536c7844c064b2e53dc0425a3eb3558732e4eb4f5e7b6d9dd1824391ab2d3a7593eb940428bf8c53a837236adaa9148891720"},
	{"85cabe17cd4d61dac94c04b180be3b5aca80229adfb77d38b0fb7b1e96dfdcc2", "f5db4adcf7fc9592af34c96c0e9855439473917426d2b905c33bc92c9d64f83185cabe17cd4d61dac94c04b180be3b5aca80229adfb77d38b0fb7b1e96dfdcc2"},
	{"0dd8f0b2be39eade722aea01f1c76e6b918eb04cef7cf55a0e64257d2a3a0c0a", "2ab407f43b2c98aefe462b49861672db1cbc9556195e612e5b3246d1778416af0dd8f0b2be39eade722aea01f1c76e6b918eb04cef7cf55a0e64257d2a3a0c0a"},
	{"d60ff5ec2999013a8464978c0b5c7fbe00ea435f53c456f3c80f3a185eef06f6", "3e273c084ac7162b16e03dc6708649fbb2a33415a1d9a81ba8ca095a95d1be21d60ff5ec2999013a8464978c0b5c7fbe00ea435f53c456f3c80f3a185eef06f6"},
	{"267543363d461f9b8b1bf10c89d2f0e92e62867889f53b308ba542e5599d2830", "d9c3dc6bc94b7ddce18e4d544568aaf866a674860f08bd9500b661035bc594e2267543363d461f9b8b1bf10c89d2f0e92e62867889f53b308ba542e5599d2830"},
	{"b15ff6b9227f98143c146b07efc25181bdd300d80bad9c35362f13bf012e28c3", "f5265b4403d9d51db621f6ae5f47fe2683da3e59c3acef2967d7b3f1d4a985edb15ff6b9227f98143c146b07efc25181bdd300d80bad9c35362f13bf012e28c3"},
	{"d62f2e721246e0e8d5a806983f43741ef12c725be57c05a49629d4d1bfab283d", "e7693a183798c5afd81d845791db119b2e633b59c5f197d5be182e272e26fb2fd62f2e721246e0e8d5a806983f43741ef12c725be57c05a49629d4d1bfab283d"},
	{"050d4c94a3c21fcddef538fde0490cb4df6c77d02a2c6cc57c9ba323a3a1d84a", "86b28ff67c294a8502225ff1fa31e15a2f91fefc26bb6227bfb71795d60a1ab3050d4c94a3c21fcddef538fde0490cb4df6c77d02a2c6cc57c9ba323a3a1d84a"},
	{"942361a92d40d9b3e4a02b4d4fa73099868fa33f85168674bddc162b57039a62", "152a28948e3db447bc75a9171e1b073dbc81ca2ff5ede9b2242c06ea9771ec42942361a92d40d9b3e4a02b4d4fa73099868fa33f85168674bddc162b57039a62"},
	{"3bcca13692fd5dd6f88e25b9e37477297e5d721232082bcfd2715896555046f5", "5e306eebe14180a2b1b7b8703e1fccfadef8e5e334704e0e696404b78d58d1a03bcca13692fd5dd6f88e25b9e37477297e5d721232082bcfd2715896555046f5"},
	{"37ff881e5fe50c38f6713ef4b8af366cd793b28564ade6fa645cb484f73244f5", "5f68bc773ecbe1ca0e03c8c585794ba1cb4ae2d1db332a50656b5dd7c1f8976137ff881e5fe50c38f6713ef4b8af366cd793b28564ade6fa645cb484f73244f5"},
	{"23804ac1e541ccda7e71f981d12f64a748e5c33fcbb675853e6b1e3d9ede3a50", "7b0efdd44e2b7931a3ce7cecb7c856abff0ba8f976b1b03852ad484e25b0da8123804ac1e541ccda7e71f981d12f64a748e5c33fcbb675853e6b1e3d9ede3a50"},
	{"f778134fc9789b3c700481e7085183fb5346f2ca8915756bbe96bd0b9d1330e2", "9b0eb7aef5c8140d4d28f409cd6d6aea7527fb8d3f78a95a3c8efac71f28f1edf778134fc9789b3c700481e7085183fb5346f2ca8915756bbe96bd0b9d1330e2"},
	{"17b3fcdfbde7fe908fa00f5bb0b58d45b403d0f138fca0ecd99084bafc7ff929", "f53fb45a4426ae66f70b4fb143ff65db3560f3bd0c9c60406f584ed5d0832e9217b3fcdfbde7fe908fa00f5bb0b58d45b403d0f138fca0ecd99084bafc7ff929"},
	{"7f5e2c020b4beafef859d872c2f6ab554c663aeb4cd4e5f9b37ad8aef5421cb1", "bc822bf2ec0b35e59e44e5c282a605ccc782e7b64a93ba5ea251ac2698dc101e7f5e2c020b4beafef859d872c2f6ab554c663aeb4cd4e5f9b37ad8aef5421cb1"},
	{"dd7edd047dc320f34e859cb96a53ec660b12758bfca47d1af9c4fac0b35d112e", "55e7a6313e34bac3cf2090765e03e742a0a934513e28a6f33568dd97431a641edd7edd047dc320f34e859cb96a53ec660b12758bfca47d1af9c4fac0b35d112e"},
	{"a97a0a8f1fcc2bd06893f581383b16e42a99658671ffe26fd3b7f8e20057b471", "3f28fbc2f9feb0e274bc2e684d41504fc4476fdbe787b46f471db5d12ac6f318a97a0a8f1fcc2bd06893f581383b16e42a99658671ffe26fd3b7f8e20057b471"},
	{"ac6aa6c1e7fecb88d1a5c54a557434b1402b8db9d3e4a5a6ece4de30d8a94e35", "39bd3217b154429ec4ee864f02263fb5629f8e1aa33634100d67d56ae5c51f2eac6aa6c1e7fecb88d1a5c54a557434b1402b8db9d3e4a5a6ece4de30d8a94e35"},
	{"367456dd9655cc0981c090e8c40055181a1af090f7c2ab472fde9338905fa33e", "4932ca97669ec34843319d8110fbcee898f8f2c54bd100fdfebc6e308e8109a5367456dd9655cc0981c090e8c40055181a1af090f7c2ab472fde9338905fa33e"},
	{"e5907a2ce97b77622d9b3cf1b8ffaedd64ab59b2b59da49eaae664f0490ea3bc", "73bbdc556596a4489e7f1ec62404cadb0be668017679b99d5a6cba244da1b10be5907a2ce97b77622d9b3cf1b8ffaedd64ab59b2b59da49eaae664f0490ea3bc"},
	{"c75f4edb7ab4a5657eb0e8cb3a6bbc4415685474f6c37d5f48a22a208d5513b5", "3be4fea79ba843f071f2e86eab0b7dfe854099847007bd9add2545b32819ecc9c75f4edb7ab4a5657eb0e8cb3a6bbc4415685474f6c37d5f48a22a208d5513b5"},
	{"8e25fb70a3c3c2a50a3edcffadbdf6675cae9fc03657a3f0fd64e605e603f419", "b589bc435e69557ce0246d3385e220df0e9f8b4a2360e42e6b029877ee47e8998e25fb70a3c3c2a50a3edcffadbdf6675cae9fc03657a3f0fd64e605e603f419"},
	{"71138f123e3f6ee2084a67308eb1179b997bff09b13bda3cd9ac9a6dc2ae6c96", "9387a256bd2e36f28d2371116cf3a17155177962d4651b364108838daecd66c271138f123e3f6ee2084a67308eb1179b997bff09b13bda3cd9ac9a6dc2ae6c96"},
	{"60d78aa0cbd80cd0bdc92ed938cc28fd4cfca5d58f5dca0d0fa333085b89fcb7", "8fafa260310d7621c2198d1c320ad07cb43ba5f314451a1b9f5343275d595afb60d78aa0cbd80cd0bdc92ed938cc28fd4cfca5d58f5dca0d0fa333085b89fcb7"},
	{"c8da241bc856870c7fd3847270cb4c33a47277389f208c8f2532c6d9facae5c6", "51135879ce3195fdf77d335fe4b2c919c9026f07306f93d8624c1c901375302ac8da241bc856870c7fd3847270cb4c33a47277389f208c8f2532c6d9facae5c6"},
	{"4bcec5487716f512e612543fab22f683a42c9f405b3e81e91c7f900c6aae7008", "4869c09ef330bbae61298794715f90dab0aad006c9053f5503f5ddb85c19cc824bcec5487716f512e612543fab22f683a42c9f405b3e81e91c7f900c6aae7008"},
	{"7881aebe3d26c9000b9c73c718a4bbdf56d3ec90feca08c0cff9acf66d4c8794", "5b7a2416ed43828a3e254ab6c6e912426dbcb8dea35cbb977a19a6f2a9942ea87881aebe3d26c9000b9c73c718a4bbdf56d3ec90feca08c0cff9acf66d4c8794"},
	{"e897aa3117f1a7523c0e9eca5fde3b9b1cb01fa4a51cdb23540efa2f1d550798", "786bed4441847bc575756ab3dbb856807facbababb31bd135402bab84db7d80ce897aa3117f1a7523c0e9eca5fde3b9b1cb01fa4a51cdb23540efa2f1d550798"},
	{"4f798a9e06acb2ac04dfd3a160357c19a65b4c1c5b052df761b9d75696e956f7", "d7a3183c82888a6826b7ebfcec3391230d9dfd58cead3c87cc185aa14dc454df4f798a9e06acb2ac04dfd3a160357c19a65b4c1c5b052df761b9d75696e956f7"},
	{"a93ff25a8c97660ac64775a01e7fc8cc623a554381888fd151ed07f5c13cb702", "620a255496dd1d9c7994e8632bb1ee1db3e715cfbb2b4ca1680e827b3997bb09a93ff25a8c97660ac64775a01e7fc8cc623a554381888fd151ed07f5c13cb702"},
	{"033a412050bf70d2ba1df0579a2b1627a06cb16f946301c1939f32d23708700e", "09a14323e3bd5b12663f2d3744dcdc84809fee62ee760eacb8767143eb4488d9033a412050bf70d2ba1df0579a2b1627a06cb16f946301c1939f32d23708700e"},
	{"7f53579c8a424ac3424fe698f86cde667426a9831fc30bbf5c31f25ee207bfa9", "a5a140fcdab2f212568a6f3211ffd86a61237606c42893bd827ef27e57c1e5e07f53579c8a424ac3424fe698f86cde667426a9831fc30bbf5c31f25ee207bfa9"},
	{"698289fa37762a55ea0f1eec2edd8558c829fab95724cc49a436915262b94e73", "a55570cf5497d8f17e77d1d568e76c5db27d170f2476c9ae900010b271216870698289fa37762a55ea0f1eec2edd8558c829fab95724cc49a436915262b94e73"},
	{"b4fa2531b7eadc17bbdce6e771ed201e0e9d58203bacbf0ba940410c6a3b7634", "9c9943bccddaca26dd29cbeeeb28a75d5d26d4bad1e8b4ec607fa13d0ea1ef13b4fa2531b7eadc17bbdce6e771ed201e0e9d58203bacbf0ba940410c6a3b7634"},
	{"f111284b4ccc62293b328b8bdea5f2302c7b7fa8fe00479724c00e0c4c384812", "2a14b15ff1cd36321686ceebcabb7c1265c803dc016a1e6884e69516846e518ef111284b4ccc62293b328b8bdea5f2302c7b7fa8fe00479724c00e0c4c384812"},
}

type Instance struct {
	Net  string // network address, with http
	Port int    // port
}

type ClientNode struct {
	Instance Instance
	Key      []string
	Balance  int
}

type ServerNode struct {
	Instances []Instance
	Key       []string
}

type Topology struct {
	Servers []ServerNode
	Clients []ClientNode
}

func writeConfig(config []byte, name string) {
	path := fmt.Sprintf("topologies/%s.json", name)
	err := ioutil.WriteFile(path, config, 0644)
	if err != nil {
		panic(err)
	}
}

// prints topology to stdout
func main() {
	writeConfig(localSimple(), "localSimple")
	writeConfig(signTest(), "signTest")
	writeConfig(localSimpleExtended(), "localSimpleExtended")
	writeConfig(localFull(), "localFull")
	writeConfig(dockerSimple(), "dockerSimple")
	writeConfig(aws(), "aws")
}

func getClientNetwork(dockerNetwork bool) string {
	if dockerNetwork {
		return "client"
	} else {
		return "localhost"
	}
}

func getServerNetwork(dockerNetwork bool) string {
	if dockerNetwork {
		return "server"
	} else {
		return "localhost"
	}
}

func config(numClients int, numServers int, numServerInstances int, dockerNetwork bool) []byte {
	topo := Topology{}

	counter := 0
	for i := 0; i < numServers; i++ {
		var instances []Instance
		for j := 0; j < numServerInstances; j++ {
			instances = append(instances, Instance{
				Net:  getServerNetwork(dockerNetwork),
				Port: 6666 + counter,
			})
			counter++
		}

		topo.Servers = append(topo.Servers, ServerNode{
			Instances: instances,
			Key:       generatedKeysServer[i],
		})
	}

	for i := 0; i < numClients; i++ {
		topo.Clients = append(topo.Clients, ClientNode{
			Instance: Instance{
				Net:  getClientNetwork(dockerNetwork),
				Port: 5555 + i,
			},
			Balance: 100,
			Key:     generatedKeysClient[i],
		})
	}

	out, _ := json.Marshal(topo)
	return out
}

// topology optimized for full local testing (includes shards)
func localFull() []byte {
	numClients := 16
	numServers := 1
	numServerInstances := 3

	return config(numClients, numServers, numServerInstances, false)
}

// topology optimized for local testing
func localSimple() []byte {
	numClients := 3
	numServers := 1
	numServerInstances := 1

	return config(numClients, numServers, numServerInstances, false)
}

// topology optimized for local testing, extended
func localSimpleExtended() []byte {
	numClients := 16
	numServers := 1
	numServerInstances := 1

	return config(numClients, numServers, numServerInstances, false)
}

// topology optimized for aws
func aws() []byte {
	numClients := 16
	numServers := 1
	numInstances := 5

	return config(numClients, numServers, numInstances, false)
}

// topology optimized for docker
func dockerSimple() []byte {
	numClients := 16
	numServers := 1
	numInstances := 1

	return config(numClients, numServers, numInstances, true)
}

// topology optimized for the sign test
func signTest() []byte {
	numClients := 64
	numServers := 64
	numInstances := 1

	return config(numClients, numServers, numInstances, true)
}
