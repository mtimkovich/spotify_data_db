select album, artist, sum(ms_played) / 60000
from songs
group by album
order by sum(ms_played) desc
limit 20;
