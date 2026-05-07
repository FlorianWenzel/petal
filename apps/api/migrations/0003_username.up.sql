alter table users add column username text;

create unique index users_username_unique on users(username);

do $$
declare
  adj text[] := array[
    'dewy','soft','gentle','quiet','sunlit','mossy','wild','sweet',
    'hushed','golden','velvet','misty','twilit','dappled','tender',
    'rosy','amber','pale','sleepy','warm','shady','still','breezy',
    'fragrant','cozy','lush','silken','starlit','creamy','meadow'
  ];
  flw text[] := array[
    'rose','lily','peony','daisy','iris','tulip','dahlia','poppy',
    'aster','lotus','jasmine','violet','orchid','marigold','zinnia',
    'magnolia','camellia','hyacinth','primrose','sunflower','bluebell',
    'foxglove','lilac','clover','fern','sage','thistle','heather',
    'crocus','freesia'
  ];
  u record;
  attempts int;
  candidate text;
begin
  for u in select id from users where username is null loop
    attempts := 0;
    loop
      candidate := adj[1 + floor(random() * array_length(adj, 1))::int]
                || '-' ||
                   flw[1 + floor(random() * array_length(flw, 1))::int]
                || '-' ||
                   (1000 + floor(random() * 9000))::int::text;
      begin
        update users set username = candidate where id = u.id;
        exit;
      exception when unique_violation then
        attempts := attempts + 1;
        if attempts > 50 then raise; end if;
      end;
    end loop;
  end loop;
end$$;

alter table users alter column username set not null;

alter table users
  add constraint users_username_chk
  check (username ~ '^[a-z0-9]([a-z0-9-]{1,22}[a-z0-9])?$');
