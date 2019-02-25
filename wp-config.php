<?php
/**
 * The base configuration for WordPress
 *
 * The wp-config.php creation script uses this file during the
 * installation. You don't have to use the web site, you can
 * copy this file to "wp-config.php" and fill in the values.
 *
 * This file contains the following configurations:
 *
 * * MySQL settings
 * * Secret keys
 * * Database table prefix
 * * ABSPATH
 *
 * @link https://codex.wordpress.org/Editing_wp-config.php
 *
 * @package WordPress
 */

// ** MySQL settings - You can get this info from your web host ** //
/** The name of the database for WordPress */
define('DB_NAME', 'wordpress_test');

/** MySQL database username */
define('DB_USER', 'root');

/** MySQL database password */
define('DB_PASSWORD', '');

/** MySQL hostname */
define('DB_HOST', 'localhost');

/** Database Charset to use in creating database tables. */
define('DB_CHARSET', 'utf8mb4');

/** The Database Collate type. Don't change this if in doubt. */
define('DB_COLLATE', '');

/**#@+
 * Authentication Unique Keys and Salts.
 *
 * Change these to different unique phrases!
 * You can generate these using the {@link https://api.wordpress.org/secret-key/1.1/salt/ WordPress.org secret-key service}
 * You can change these at any point in time to invalidate all existing cookies. This will force all users to have to log in again.
 *
 * @since 2.6.0
 */
define('AUTH_KEY',         'Yp)+K76fqFAz=.eZSy<N5.xK|worV6*[6=oU>N8j7~Mbe&h<5m{f%YKr$OoA]A>A');
define('SECURE_AUTH_KEY',  'R/Diy#]LXzaAI3=g2(A[762T-` %ujWr{G36Dl9aM[34,rI!DAPZkzY;d?)INkFj');
define('LOGGED_IN_KEY',    'U{IT$`<=;qu(&RB?>GZThtIQpET/{d{ 2V KkQ.O-#L&nfYZUs_QKO#UH2|f &GW');
define('NONCE_KEY',        '!,:Z)K3239|gE]@GP|u-2pq5YR,/6(=WEOFAM REyZ129yW> 568*N1Xd#$1VE#P');
define('AUTH_SALT',        'DoFBX),s:{uWN*mX>rM*F$ V&({QI8mXg<o}j`F!C@(P$XB[MV>;&@vESK|8;@$z');
define('SECURE_AUTH_SALT', '3bfe}_y>ia@X=mYgZ4w%d?`e,Er)RT`&{_$syyQv?7 0w,$L#Ggi@1O_Z }S>O<u');
define('LOGGED_IN_SALT',   '1QPVr>LjwOTHi&fk+c<~DZLDpBNv1xTr0BMth6k,-6^>EfztXcHN>~B0YsDU%jss');
define('NONCE_SALT',       '6BuA*mriwX3vC_a4Tb2DH$(bo/<(e V3W2.H{}g]zq]*b{WQPDu0[M:)a>=QHEfz');

/**#@-*/

/**
 * WordPress Database Table prefix.
 *
 * You can have multiple installations in one database if you give each
 * a unique prefix. Only numbers, letters, and underscores please!
 */
$table_prefix  = 'wp_';

/**
 * For developers: WordPress debugging mode.
 *
 * Change this to true to enable the display of notices during development.
 * It is strongly recommended that plugin and theme developers use WP_DEBUG
 * in their development environments.
 *
 * For information on other constants that can be used for debugging,
 * visit the Codex.
 *
 * @link https://codex.wordpress.org/Debugging_in_WordPress
 */
define('WP_DEBUG', false);

/* That's all, stop editing! Happy blogging. */

/** Absolute path to the WordPress directory. */
if ( !defined('ABSPATH') )
	define('ABSPATH', dirname(__FILE__) . '/');

/** Sets up WordPress vars and included files. */
require_once(ABSPATH . 'wp-settings.php');
