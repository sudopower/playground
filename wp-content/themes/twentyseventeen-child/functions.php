<?php
add_action( 'wp_enqueue_scripts', 'my_theme_enqueue_styles' );
function my_theme_enqueue_styles() {
 
    $parent_style = 'twentyseventeen-style'; 
 
    wp_enqueue_style( $parent_style, get_template_directory_uri() . '/style.css' );
    wp_enqueue_style( 'child-style',
        get_stylesheet_directory_uri() . '/style.css',
        array( $parent_style ),
        wp_get_theme()->get('Version')
    );
}

function add_mix_it_up() {
	wp_enqueue_style( 'mix-it-up', get_stylesheet_directory_uri() . '/assets/css/mix_it_up.css'  );
	wp_enqueue_script( 'mix-it-up', get_stylesheet_directory_uri() . '/assets/js/mixitup.min.js',$deps = array(), $ver = false, $in_footer = true);
	wp_enqueue_script( 'my-custom-scripts', get_stylesheet_directory_uri() . '/assets/js/scripts.js',$deps = array(), $ver = false, $in_footer = true);
}

add_action( 'wp_enqueue_scripts', 'add_mix_it_up');

function get_post_gallery_images_custom( $content ) {

    global $post;

    // Only do this on singular items
 	if( ! is_singular() )
     return $content;

 // Make sure the post has a gallery in it
 if( ! has_shortcode( $post->post_content, 'gallery' ) )
     return $content;

 // Retrieve the first gallery in the post
 $gallery = get_post_gallery_images( $post );

$image_list = '<ul>';

// Loop through each image in each gallery
foreach( $gallery as $image_url ) {

    $image_list .= '<li>' . '<img src="' . $image_url . '">' . '</li>';

}

$image_list .= '</ul>';

// Append our image list to the content of our post
$content .= $image_list;

 return $content;

}
add_filter( 'the_content', 'get_post_gallery_images_custom' );
?>