 # MLD — Modèle Logique de Données
                                                                                                                                                                                                                                             
  Users(                                                                                                                                                                                                                                   
    user_id       PK,
    username      UNIQUE NOT NULL,
    email         UNIQUE NOT NULL,                                                                                                                                                                                                           
    password      NOT NULL,           -- bcrypt hash; NULL if OAuth-only account
    avatar_url,                                                                                                                                                                                                                              
    role          NOT NULL,           -- member | moderator | admin                                                                                                                                                                        
    created_at, updated_at, deleted_at                                                                                                                                                                                                       
  )                                                                                                                                                                                                                                        
                                                                                                                                                                                                                                             
  Posts(                                                                                                                                                                                                                                   
    post_id       PK,
    author_id     FK → Users(user_id) NOT NULL,
    title         NOT NULL,                                                                                                                                                                                                                  
    content       NOT NULL,           -- HTML/JSON from Tiptap
    status        NOT NULL,           -- pending_moderation | approved | flagged | blocked                                                                                                                                                   
    image_url,                                                                                                                                                                                                                             
    embedding     vector(768),        -- nomic-embed-text via Ollama + pgvector                                                                                                                                                              
    created_at, updated_at, deleted_at                                                                                                                                                                                                       
  )                                                                                                                                                                                                                                          
                                                                                                                                                                                                                                             
  Categories(                                                                                                                                                                                                                              
    category_id   PK,
    name          UNIQUE NOT NULL,
    created_at, updated_at
  )

  PostCategories(                                                                                                                                                                                                                            
    post_id       FK → Posts(post_id),
    category_id   FK → Categories(category_id),                                                                                                                                                                                              
    PRIMARY KEY (post_id, category_id)                                                                                                                                                                                                     
  )

  Comments(
    comment_id    PK,
    author_id     FK → Users(user_id) NOT NULL,                                                                                                                                                                                              
    post_id       FK → Posts(post_id) NOT NULL,
    content       NOT NULL,                                                                                                                                                                                                                  
    status        NOT NULL,           -- same enum as Posts                                                                                                                                                                                
    created_at, updated_at, deleted_at                                                                                                                                                                                                       
  )
                                                                                                                                                                                                                                             
  PostLikes(                                                                                                                                                                                                                               
    post_like_id  PK,
    user_id       FK → Users(user_id) NOT NULL,
    post_id       FK → Posts(post_id) NOT NULL,                                                                                                                                                                                              
    liked         BOOLEAN NOT NULL,
    created_at,                                                                                                                                                                                                                              
    UNIQUE (user_id, post_id)                                                                                                                                                                                                              
  )                                                                                                                                                                                                                                          
   
  CommentLikes(                                                                                                                                                                                                                              
    comment_like_id  PK,                                                                                                                                                                                                                   
    user_id          FK → Users(user_id) NOT NULL,
    comment_id       FK → Comments(comment_id) NOT NULL,
    liked            BOOLEAN NOT NULL,                                                                                                                                                                                                       
    created_at,
    UNIQUE (user_id, comment_id)                                                                                                                                                                                                             
  )                                                                                                                                                                                                                                        

  OAuthProviders(
    provider_id       PK,
    user_id           FK → Users(user_id) NOT NULL,                                                                                                                                                                                          
    provider_name     NOT NULL,       -- github | google
    provider_user_id  NOT NULL,                                                                                                                                                                                                              
    created_at, updated_at,                                                                                                                                                                                                                
    UNIQUE (provider_name, provider_user_id)                                                                                                                                                                                                 
  )                                                                                                                                                                                                                                          
   
  RefreshTokens(                                                                                                                                                                                                                             
    token_id      PK,                                                                                                                                                                                                                      
    user_id       FK → Users(user_id) NOT NULL,
    token_string  UNIQUE NOT NULL,
    expires_at    NOT NULL,
    created_at                                                                                                                                                                                                                               
  )
                                                                                                                                                                                                                                             
  PasswordResetTokens(                                                                                                                                                                                                                     
    token_id      PK,
    user_id       FK → Users(user_id) NOT NULL,
    token_string  UNIQUE NOT NULL,
    expires_at    NOT NULL,                                                                                                                                                                                                                  
    used          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at                                                                                                                                                                                                                               
  )                                                                                                                                                                                                                                        

  Notifications(
    notification_id  PK,
    user_id          FK → Users(user_id) NOT NULL,
    type             NOT NULL,   -- new_comment | new_like | new_reply | new_dm | post_flagged | post_approved                                                                                                                               
    title            NOT NULL,                                                                                                                                                                                                               
    content          NOT NULL,                                                                                                                                                                                                               
    is_read          BOOLEAN NOT NULL DEFAULT FALSE,                                                                                                                                                                                         
    created_at                                                                                                                                                                                                                               
  )                                                                                                                                                                                                                                        

  Messages(
    message_id   PK,
    sender_id    FK → Users(user_id) NOT NULL,
    receiver_id  FK → Users(user_id) NOT NULL,                                                                                                                                                                                               
    content      NOT NULL,
    is_read      BOOLEAN NOT NULL DEFAULT FALSE,                                                                                                                                                                                             
    created_at                                                                                                                                                                                                                             
  )

  Reports(
    report_id     PK,
    reporter_id   FK → Users(user_id) NOT NULL,                                                                                                                                                                                              
    content_type  NOT NULL,   -- post | comment
    content_id    NOT NULL,   -- ID of the reported post or comment                                                                                                                                                                          
    reason        NOT NULL,                                                                                                                                                                                                                  
    status        NOT NULL DEFAULT 'pending',  -- pending | resolved | dismissed
    created_at                                                                                                                                                                                                                               
  )                                                                                                                                                                                                                                        
                                                                                                                                                                                                                                             
  PushSubscriptions(                                                                                                                                                                                                                       
    subscription_id  PK,
    user_id          FK → Users(user_id) NOT NULL,
    endpoint         NOT NULL,                                                                                                                                                                                                               
    p256dh_key       NOT NULL,
    auth_key         NOT NULL,                                                                                                                                                                                                               
    created_at                                                                                                                                                                                                                             
  )
