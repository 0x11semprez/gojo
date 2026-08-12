why gojo ? </br>
coldcard was hacked due to a conditional/build configuration that caused the secure RNG path not to be used as intended.</br>
</br>
<https://www.trmlabs.com/resources/blog/the-largest-hardware-wallet-exploit-of-2026-inside-the-usd-116-million-coldcard-hack> </br>
</br>
gojo is not a wallet. </br>
gojo is a custodial KMS (key management service) for bitcoin and monero keys, exposed through a golang api. </br>
gojo generates extremely secure keys by mixing several hardware-specific entropy sources in rust, instead of trusting a single RNG path. </br>
gojo is gojo </br>
</br>
architecture: </br>

- postgres: keeps the encrypted keys and the accounts </br>
- golang api: handles accounts and authentication, stores/serves keys, always encrypted </br>
- rust: generates the keys (bitcoin secp256k1, monero ed25519), collecting multiple hardware entropy sources before seeding </br>
  </br>
  no need to have a long readme, just check all schemas in the docs files. we focused more on design system because coding skills are not as valuable in the eyes of others anymore... </br>
  </br>
  we love coding
