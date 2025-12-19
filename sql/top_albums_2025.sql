select album, artist, sum(ms_played) / 60000
from songs
where strftime('%Y', timestamp) = '2025'
group by album
order by sum(ms_played) desc
limit 20;
