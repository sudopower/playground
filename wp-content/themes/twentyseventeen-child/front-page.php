<?php
get_header(); ?>

<div class="wrap">
	<div id="primary" class="content-area">
             <!-- POST LOOP START -->
        <div class="controls">
            <button type="button" class="filter" data-filter="all">All</button>
            <?php		
            global $post;
			$showcasePosts = new WP_Query('cat=7');			
			while($showcasePosts->have_posts()) :   $showcasePosts->the_post();	
			?>
            <!-- LOOP -->                
            <button type="button" class="filter" data-filter=".<?php echo $post->post_name;?>"><?php the_title();?></button>
			<?php endwhile;			
			//wp_reset_postdata();
			?>
        </div>
		 <!-- POST LOOP END -->
        <div class="showcase-gallery">
            <?php	          	  
                while($showcasePosts->have_posts()) :   $showcasePosts->the_post();	
                ?>
                <!-- LOOP -->                
                
                <?php 
                    $gallery = get_post_gallery_images($post = 66);   
                    print_r($gallery);                 
                    foreach( $gallery as $image_url ) {?>
                        <a href="<?php the_permalink();?>" class="mix<?php echo ' '.$post->post_name;?>">
                            <img src="<?php $image_url ?>">                
                        </a>
                <?php }
                endwhile;			
                wp_reset_postdata();
                ?>
        </div>
	</div><!-- #primary -->        
</div><!-- .wrap -->

<?php
get_footer();
