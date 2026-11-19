// file generated. see /config.
package config

import "encoding/json"

var Config WeshConfig

// FIXME: make it more nicely
func init() {
	input := `
{
  "wesh": {
    "contacts": {
      "testbot": {
        "link": "https://tone.onl/id#contact/oZBLExQaU8TTubcawgfs8AVcH9FbFV5J1kdEVqgueBaFCMWsADV341pCv26zcHMv3GaTL6UHPzVoHzLoereFQLmBGAh2QKb/name=TestBot",
        "name": "TestBot",
        "description": "Official TestBot",
        "kind": "Bot",
        "suggestion": false,
        "icon": ""
      },
      "testbot-dev": {
        "link": "https://tone.onl/id#contact/oZBLFRW7b1CthoZojwAJTSv27Pg7XAeqHKTDnpCkhz8wx2vGNxef4L6X1Eyh2jB9QdhWFtC6WNpEq6S4nGhZMCfFYEacZPS/name=TestBot-Dev",
        "name": "TestBot Dev",
        "description": "TestBot-Dev",
        "kind": "Bot",
        "suggestion": false,
        "icon": ""
      },
      "welcomebot": {
        "link": "https://tone.onl/id#contact/oZBLFukkeefDcXibkvxHyd1kZpqUUm4mKsE3tgcDakWx6Pa5Ak3VghVyxLZ8vns71EDLWkLbCZ5QcGADnehvxyq3B7LFh14/name=WelcomeBot",
        "name": "Welcome Bot",
        "description": "Official Welcome Bot",
        "kind": "Bot",
        "suggestion": true,
        "icon": "wesh_bot_pink_bg"
      },
      "welcomebot-dev": {
        "link": "https://tone.onl/id#contact/oZBLFQQTqmZ25yvcRRUVgzzZCpDY7YAVXDEz3g8kmwcSPCmnQrqiMbLrfSp1GRLaNBX6ZRarid8jkSUX4yQYL4oBZcLDtFY/name=WelcomeBot-Dev",
        "name": "Welcome Dev",
        "description": "Welcome Bot - Dev",
        "kind": "Bot",
        "suggestion": false,
        "icon": "wesh_dev_green_bg"
      }
    },
    "conversations": {
      "dev-test": {
        "link": "https://tone.onl/id#group/5QdUv6Fn3uvi3AchNcqFECaZvngDzWGUJV4wntYCWuYyjXrNi4ykvnP2ZCeWA1YLTmCSFbZaimXzp8rZtkKqby8nGv7S2AXyJPggo3aMghJS3rpeuDx6pbbNqEoNXJK67pWVpcLB6VrXM54zRNbi4zhEsLdATGzguWoruShcRHzFpuP75nyvRnXnfRaHPjdbyC/name=Dev+Test",
        "name": "Dev Test",
        "description": "Simple Conversation with Replication",
        "suggestion": false
      },
      "dev-test-2": {
        "link": "https://tone.onl/id#group/5QdUv6Fn3uwSPQZpLUXcNpQ24bR9y1TUK6xJfE9khJo4jrP8jF4QfdgpTT73Us58y6XZhtvJHLDCouXfgsKDxXqZAsfLDhUZLb48PDcBYYVUk7nzxh1MKh5a6Ug1bP4KpdrVDYQx1gCZA4HZavEUgUBbC1pYBZ2DaY27sCMfqxpt79RjZCBwmVqW1DbbuYCUoc/name=Dev+Test+2",
        "name": "Dev Test 2",
        "description": "Simple Conversation with Replication",
        "suggestion": false
      },
      "group-with-passphrase": {
        "link": "https://tone.onl/id#enc/3KE388cyt693i3Y8FUK1HMwr6BQGuckRx8ypxdtmrFrmMvgb8JNeu8L8BSkL6MbNjyFR4jSgEDmerbmGSXPsoVrbmrprcKX3PG3a2EY4jF2D7BUUdg1kDpxewAcHMyVRW66yMTe29G6sjuJ8C2MjEUWMWT7oX3jv94aCvWQbomePYL6ff9Bm8Msbe3UXBgqN6Wff9QqEScsMG9WnKUQQLtKxSqwN59etjrz/name=group-with-passphrase",
        "name": "Group with passphrase (test)",
        "description": "A demo group with a passphrase",
        "suggestion": false
      }
    }
  },
  "p2p": {
    "rdvp": [
      {
        "maddr": "/ip4/51.15.25.224/udp/4040/quic-v1/p2p/12D3KooWHhDBv6DJJ4XDWjzEXq6sVNEs6VuxsV1WyBBEhPENHzcZ"
      }
    ],
    "static-relays": [
      "/ip4/51.15.25.224/udp/6363/quic-v1/p2p/12D3KooWAHcEz4K5XAgRDav9fLuhiRY2wuXip385EmT5RoRkCmjr"
    ]
  }
}`
	err := json.Unmarshal([]byte(input), &Config)
	if err != nil {
		panic(err)
	}
}
