# Shortcuts

| Key | Action |
| --- | --- |
|`F2`|Import AlfaBank account statements (.pdf)|
|`F3`|Page with statistics by selected account and date range|
|`Esc`|Close page|
|`Ctrl + A`|- Insert into **Transactions** table,<br>- Pick file in **Attachments** table<br>- Add account on focused **Accounts** box<br>- Add category on focused **Categories** box|
|`Ctrl + D`|- Delete selected row in **Transactions** table,<br>- Delete file in **Attachments** table<br>- Delete selected account on focused **Accounts** box<br>- Delete selected category on focused **Categories** box|
|`Ctrl + R`|- Edit selected account on focused **Accounts** box<br>- Edit selected category on focused **Categories** box|
|`Enter`|- Edit selected transaction on focused **Transactions** box<br>- Open selected file with default system application on focused **Attachments** box|
|`Ctrl + C`|Exit|

>[!note]
>- Attachments files are only saved after adding a new transaction or saving existing one.
>- On first open application creates configuration folder with `tree.log` and `config.yml`:
>	- Linux - `$HOME/.config/money-tree`
>	- Windows - `%APPDATA%\money-tree`
>- Also in directory where application was started it creates `attachemnts` folder and `database.db` file.

# `config.yml`:
```yaml
path: d:\test\moneytree
database: d:\test\moneytree\database.db
attachments: d:\test\moneytree\attachments
```
