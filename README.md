# Balance

[![Hits-of-Code][1]][2]

[1]: https://hitsofcode.com/github/slotopol/balance?branch=main
[2]: https://hitsofcode.com/github/slotopol/balance/view?branch=main

GUI tool application that helps to change users balance for Slotopol server.

On application starting, you need to log in with an account that has administrator rights to work with users. Next, you can make a list of accounts to work with. This list is automatically saved. The buttons on the toolbar allow you to top up the balance for the selected user by specifying the amount of replenishment. You can change the RTP of the selected user, in all games the reels will be selected with payments closest to the specified percentage. You can also change the access rights for accounts if your account itself has the appropriate rights. If your account has the right to access the club management, you can replenish / withdraw money from the bank, jackpot fund, and deposit.

Screenshot:

[![balance #1](assets/scrn-icon.webp)](assets/screenshot.webp)

When you start a new server instance with an empty database, 3 new accounts are automatically registered: admin, dealer, and player. You need to log in on start with the admin account `admin@example.org` with the password `0YBoaT`, after which you can perform any manipulations.

Note! Since the server is designed for monopoly usage of the database, direct editing of the database is only possible when the server is not running.

## About MRTP

MRTP - Master Return to Player percentage. Each game in most common cases have the set of reels with different RTP. Before any spin in any game at first takes the reel with RTP closest to pointed value of MRTP. RTP of reel can be more of MRTP or less of it, but its closest then others. MRTP can be in range from 85 to 115. Original reels in most common cases have RTP near 95%, but this percetage can not be changed at source providers. This server provides reels for most of games with step of RTP 1-2% from 85% to 100% and some reels with unprofitable RTP. Reels with RTP > 100% can be used in demo games. Some games, such as `Sizzling Hot` have a very high pays at paytable, and the reels with low RTP cannot be composed in a way that makes the game looks good. You can explore the list of available reels for each game by command

```sh
slot_win_x64.exe list --all --rtp
```

## About access rights

At right in application interface you can see `access` column with users access rights. Access rights divided by 5 types, which can be combined:

* **member** - account have access to club. This access must have any account to play the games. You can ban a user by removing this access.

* **dealer** - account can change club game settings and users gameplay. With this access account can play any games in club instead of another user. It's for game dealers.

* **booker** - account can change user properties and move user money to/from club deposit. This access must have club cashiers.

* **master** - account can change club bank, fund, deposit. This access points to club master.

* **admin** - account can change same access levels to other users. For server administrators.

---
(c) schwarzlichtbezirk, 2024.
