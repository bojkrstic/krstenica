# krstenica


1. Projekat Krstenica sadrzi api koji se pokrece preko docker compose
- docker compose up -d   -- na ovaj nacin se pokrene dokerizovani posrgesSql
- docker compose down    -- obaranje istog dockera
Ovde je vazno napomenuti da je posgresSQL dokerizovan, i da se sve radi u docker kontejneru
Citanje verzije postgresSql
- psql --version

2. Da li uspostavljena veza sa postgreSQL-om moze da se uradi na ovaj nacin
- psql "postgres://postgres:bokana@localhost:5432/hrams?sslmode=disable"   -- to se pokrece iz terminala

3. Zatim mora da se pokrene sam api ali mora d a se nalazis u folderu gde je main.go, a to je folder /home/krle/develop/horisen/Krstenica-new/krstenica/cmd
- ./krstenica-api --config-file-path=/home/krle/develop/horisen/Krstenica-new/krstenica/config/krstenica_api_conf.json
ili skraceno ako se nalazimo u folderu vec
- ./krstenica-api --config-file-path=../config/krstenica_api_conf.json

4. Neke vazne komande za rad sa dockerom za postgresSQL
- docker ps  -- provera da li je podignut kontejner
- docker exec -it krstenica-db sh --ovako se ulazi u sami docker file, i tu se dalje radi sa postresSql direktno u kontjeneru
- psql -U postgres  -- prvo moras da udjes u postgres
- \l  -- ako zelis da vidis sve baze
- \dt -- ako zelis da vidis sve tabele
- \c ime_baze
- postgres=# \c hrams  -- ovo je odziv postgres=#
- hrams=#   -- sad je odziv posle ove komande \c hrams
- INSERT INTO hram (hram_id, naziv_hrama, created_at) VALUES (1, 'Presveta Bogorodica', '2024-11-30 22:55:00+01');  -- unos
- SELECT hram_id FROM public.hram;  -- list
- \d nazivtabele -- proverava strukturu tabele
- Automatsko generisanje vrednosti: Ako želite da PostgreSQL automatski generiše vrednosti za hram_id, možete postaviti hram_id kolonu da koristi sekvencu:
- CREATE SEQUENCE hram_id_seq;
ALTER TABLE public.hram ALTER COLUMN hram_id SET DEFAULT nextval('hram_id_seq');

- INSERT INTO public.hram (naziv_hrama, created_at) VALUES ('Presveta Bogorodica', '2024-11-30 22:55:00+01');
- SELECT * FROM public.hram;
------------------------------------------------------------------------------------------------------------------------------------------------
postgres=# \c
You are now connected to database "postgres" as user "postgres".
postgres=# \l
postgres=# \c hrams
You are now connected to database "hrams" as user "postgres".
hrams=# SELECT * FROM public.hram;
 hram_id |           naziv_hrama            |          created_at           
---------+----------------------------------+-------------------------------
       1 | Presveta Bogorodica              | 2024-11-30 21:55:00+00
       2 | Sveti Arhangel Mihailo           | 2024-11-30 21:56:00+00
       3 | Presveta Bogorodica Ciniglavci   | 2024-11-30 21:55:00+00
       4 | 0                                | 2024-11-30 22:52:47.681376+00
      14 | Bojan 3812638179211989351        | 2024-11-30 23:37:13.825635+00
      15 | Bojan Krstic 3812638179211989351 | 2024-11-30 23:38:25.602259+00
(6 rows)

hrams=# \d hram
                                       Table "public.hram"
   Column    |           Type           | Collation | Nullable |             Default              
-------------+--------------------------+-----------+----------+----------------------------------
 hram_id     | integer                  |           | not null | nextval('hram_id_seq'::regclass)
 naziv_hrama | character varying(100)   |           | not null | 
 created_at  | timestamp with time zone |           | not null | CURRENT_TIMESTAMP
Indexes:
    "hram_pkey" PRIMARY KEY, btree (hram_id)
    "hram_naziv_hrama_key" UNIQUE CONSTRAINT, btree (naziv_hrama)

------------------------------------------------------------------------------------------------------------------------------------------------------
- \q -- exit,logout


5. Potencijani problemi koji su se javljali
- 1. _ "github.com/lib/pq"   -- ovo je drajver za postgresSQL, obavezan je i treba da se nalazi u dao delu
- 2. zatim obavezno je da postoji port u docker compose file-u, deo ispod
    -  ports:
      - "5432:5432"
- 3. password mora da se poklopi izmedju docker compose file-a i samog json fajla za konfiguraciju.


6. Da bi napravio izvrsnu verziju apija u folderu gde se nalazi main.go pokrne se
- go build -o krstenica-api
napravi se krstenica-api.go izvrsna datoka binarna
7. Ovaj deo vec imas u 3. 
Zatim mora da se pokrene sam api ali mora d a se nalazis u folderu gde je main.go, a to je folder /home/krle/develop/horisen/Krstenica-new/krstenica/cmd
- ./krstenica-api --config-file-path=/home/krle/develop/horisen/Krstenica-new/krstenica/config/krstenica_api_conf.json
ili skraceno ako se nalazimo u folderu vec
- ./krstenica-api --config-file-path=../config/krstenica_api_conf.json



8. Problem sa regularnim izrazom u api_path_registry.go
Naime problem je bio jer nisam mogao dobro da preuzmem parametar, Problem se resava na dva nacina.
- 1. Prvi nacin je da se umesto literala `` ubaci "" i onda je sve ok, jer se u tom slucaju prepozna regularni izraz i $id
- 2. Drugi nacin da se promeni sam regularni izraz sa: 
        re := regexp.MustCompile(`\$[a-zA-Z0-9]+`) // Ispravan regex za prepoznavanje dinamičkih parametara
Ovo je pravilo problem: 
        re := regexp.MustCompile(`\\$[a-zA-Z9-9]+`)
        re := regexp.MustCompile("\\$[a-zA-Z9-9]+")  Ovo je ispravan isto da ide sa "".
        

9. Citanje logova
- docker logs krstenica-db
10. Citanje volumena
- docker volume ls
11. Brisanje volumena
- docker-compose down --volumes
12. Dodavanje polja u tableli
- ALTER TABLE hram ADD COLUMN IF NOT EXISTS status VARCHAR(50);
to se dodaje u init.sql fajlu(recimo init-hram-db.sql)
13. Ovako se dodaje ako se koristi init.sql
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hram') THEN
        -- Dodaj kolonu 'status' ako ne postoji
        ALTER TABLE hram ADD COLUMN IF NOT EXISTS status VARCHAR(50);
    END IF;
END $$;
