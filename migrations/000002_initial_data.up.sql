INSERT INTO users (username, password, is_admin)
VALUES (
        'admin',
        '$2a$10$aOFCGN3YY8TWxj6WLSfTceuCb0LkKEnuuhbj3adAH7ET7UAwxhIja',
        true
    ),
    (
        'user',
        '$2a$10$xoSLkRdTwcOSNHoqULaiHukk04iNTBTlgEO203GJY6qm0jQWQKb9e',
        false
    );
INSERT INTO actors (
        id,
        first_name,
        last_name,
        middle_name,
        sex,
        birthday
    )
VALUES (
        1,
        'John',
        'Doe',
        'J',
        'male',
        '2006-01-02 00:00:00'
    ),
    (
        2,
        'Ryan',
        'Gosling',
        NULL,
        'male',
        '1980-11-12 00:00:00'
    ),
    (
        3,
        'Ryan',
        'Reynolds',
        NULL,
        'male',
        '1976-10-23 00:00:00'
    ),
    (
        4,
        'Margot',
        'Robbie',
        NULL,
        'female',
        '1990-07-02 00:00:00'
    );
INSERT INTO films (id, title, description, release_date, rating)
VALUES (
        1,
        'Drive',
        'A stuntman and getaway driver falls in love with Irene who is married to a criminal. In a bid to protect her from her husband and some gangsters, he decides to cross over to the other side of the law.',
        '2011-11-03',
        7
    ),
    (
        2,
        'Blade Runner 2049',
        'K, an officer with the Los Angeles Police Department, unearths a secret that could create chaos. He goes in search of a former blade runner who has been missing for over three decades.',
        '2017-10-05',
        8
    ),
    (
        3,
        'Barbie',
        'Barbie and Ken are having the time of their lives in the colorful and seemingly perfect world of Barbie Land. However, when they get a chance to go to the real world, they soon discover the joys and perils of living among humans.',
        '2023-07-20',
        7
    ),
    (
        4,
        'The Proposal',
        'When New York editor Margaret faces deportation, she convinces her assistant Andrew to marry her in return for a promotion. However, when she visits his hometown, it changes her in many ways.',
        '2009-06-18',
        8
    ),
    (
        5,
        'The Wolf of Wall Street',
        'Introduced to life in the fast lane through stockbroking, Jordan Belfort takes a hit after a Wall Street crash. He teams up with Donnie Azoff, cheating his way to the top as his relationships slide.',
        '2013-02-06',
        9
    );

INSERT INTO actors_and_films (actor_id, film_id)
VALUES (2, 1), (2,2), (2,3), (4,3), (3,4), (4,5)