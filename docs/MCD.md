
  # MCD — Modèle Conceptuel de Données
                                                                                                                                                                                                                                             
  ## Entities                                                                                                                                                                                                                                 
                                                                                                                                                                                                                                             
  Users(user_id, username, avatar_url, email, password, role, created_at)                                                                                                                                                                  

  Posts(post_id, author_id, title, content, status, image_url, embedding, created_at, updated_at)                                                                                                                                            
   
  Categories(category_id, name)                                                                                                                                                                                                              
                                                                                                                                                                                                                                           
  PostCategories(post_id, category_id)                                                                                                                                                                                                       
   
  Comments(comment_id, author_id, post_id, status, content, created_at, updated_at)                                                                                                                                                          
                                                                                                                                                                                                                                           
  PostLikes(post_like_id, user_id, post_id, liked)                                                                                                                                                                                           
   
  CommentLikes(comment_like_id, user_id, comment_id, liked)                                                                                                                                                                                  
                                                                                                                                                                                                                                           
  OAuthProviders(provider_id, user_id, provider_name, provider_user_id)                                                                                                                                                                      
   
  RefreshTokens(token_id, user_id, token_string, expiration_date)                                                                                                                                                                            
                                                                                                                                                                                                                                           
  PasswordResetTokens(token_id, user_id, token_string, expiration_date, used)                                                                                                                                                                
   
  Notifications(notification_id, user_id, type, title, content, is_read, created_at)                                                                                                                                                         
                                                                                                                                                                                                                                           
  Messages(message_id, sender_id, receiver_id, content, is_read, created_at)                                                                                                                                                                 
   
  Reports(report_id, reporter_id, content_type, content_id, reason, status, created_at)                                                                                                                                                      
                                                                                                                                                                                                                                           
  PushSubscriptions(subscription_id, user_id, endpoint, p256dh_key, auth_key, created_at)                                                                                                                                                    
   
  ## Associations                                                                                                                                                                                                                               
                                                                                                                                                                                                                                           
  A user can write several posts                    → (0,n)
  A post is written by a single user                → (1,1)                                                                                                                                                                                  
                                                                                                                                                                                                                                             
  A post can belong to several categories           → (1,n)                                                                                                                                                                                  
  A category can be associated with several posts   → (0,n)                                                                                                                                                                                  
                                                                                                                                                                                                                                             
  A post can have several PostLikes                 → (0,n)                                                                                                                                                                                  
  A PostLike belongs to a single post               → (1,1)
  A user can have several PostLikes                 → (0,n)                                                                                                                                                                                  
  A PostLike is associated with a single user       → (1,1)                                                                                                                                                                                  
                                                                                                                                                                                                                                             
  A post can have several comments                  → (0,n)                                                                                                                                                                                  
  A comment belongs to a single post                → (1,1)                                                                                                                                                                                  
  A user can write several comments                 → (0,n)                                                                                                                                                                                
  A comment is written by a single user             → (1,1)                                                                                                                                                                                  
   
  A comment can have several CommentLikes           → (0,n)                                                                                                                                                                                  
  A CommentLike belongs to a single comment         → (1,1)                                                                                                                                                                                
  A user can have several CommentLikes              → (0,n)                                                                                                                                                                                  
  A CommentLike is associated with a single user    → (1,1)                                                                                                                                                                                
                                                                                                                                                                                                                                             
  A user can be linked to several OAuthProviders    → (0,n)
  An OAuthProvider is linked to a single user       → (1,1)                                                                                                                                                                                  
                                                                                                                                                                                                                                             
  A user can have several RefreshTokens             → (0,n)
  A RefreshToken is linked to a single user         → (1,1)                                                                                                                                                                                  
                                                                                                                                                                                                                                             
  A user can have several PasswordResetTokens       → (0,n)
  A PasswordResetToken is linked to a single user   → (1,1)                                                                                                                                                                                  
                                                                                                                                                                                                                                           
  A user can receive several Notifications          → (0,n)                                                                                                                                                                                  
  A Notification is linked to a single user         → (1,1)
                                                                                                                                                                                                                                             
  A user can send several Messages                  → (0,n)                                                                                                                                                                                
  A user can receive several Messages               → (0,n)
  A Message has a single sender                     → (1,1)                                                                                                                                                                                  
  A Message has a single receiver                   → (1,1)
                                                                                                                                                                                                                                             
  A user can file several Reports                   → (0,n)                                                                                                                                                                                  
  A Report is filed by a single user                → (1,1)
                                                                                                                                                                                                                                             
  A user can have several PushSubscriptions         → (0,n)                                                                                                                                                                                
  A PushSubscription belongs to a single user       → (1,1)
                                                                                                                                                                                                                                             
 ## Cardinalities
                                                                                                                                                                                                                                             
  Users       (0,n)  --writes-->   (1,1)  Posts                                                                                                                                                                                            
  Posts       (1,n)  --linked-->   (0,n)  Categories                                                                                                                                                                                         
  Posts       (0,n)  --has-->      (1,1)  PostLikes
  Posts       (0,n)  --has-->      (1,1)  Comments                                                                                                                                                                                           
  Comments    (0,n)  --has-->      (1,1)  CommentLikes                                                                                                                                                                                     
  Users       (0,n)  --writes-->   (1,1)  Comments                                                                                                                                                                                           
  Users       (0,n)  --has-->      (1,1)  PostLikes                                                                                                                                                                                          
  Users       (0,n)  --has-->      (1,1)  CommentLikes                                                                                                                                                                                       
  Users       (0,n)  --linked-->   (1,1)  OAuthProviders                                                                                                                                                                                     
  Users       (0,n)  --linked-->   (1,1)  RefreshTokens                                                                                                                                                                                    
  Users       (0,n)  --linked-->   (1,1)  PasswordResetTokens                                                                                                                                                                                
  Users       (0,n)  --receives--> (1,1)  Notifications                                                                                                                                                                                    
  Users       (0,n)  --sends-->    (1,1)  Messages                                                                                                                                                                                           
  Users       (0,n)  --receives--> (1,1)  Messages                                                                                                                                                                                         
  Users       (0,n)  --files-->    (1,1)  Reports                                                                                                                                                                                            
  Users       (0,n)  --has-->      (1,1)  PushSubscriptions 

