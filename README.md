# pock_multiplayer
pock for a multiplayer game


Gateway/Gateway

Route :
  - login
  - logout
  - refresh
  - reset_password
  - register
  - user_world
  - world_list
  - join_world
  - delete_user_world


NoSQL Database :
  - id_world : 
    - id_user :
      - name
      - posXY
      - customization :
        - hair
        - nose
        - ...
      - stats
      - equipment :
        - weapon
        - armor
        - ...
      - inventory
      - job

SQL Database : 
    - user
    - user_world_aff
    - connexion
    - device
    - world_list