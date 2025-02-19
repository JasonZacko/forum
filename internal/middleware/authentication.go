// middleware/authentication.go
package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"
)

func Authentication(db *sql.DB, cache *SessionCache, protectedPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the requested path requires authentication
			requiresAuth := false
			for _, path := range protectedPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					requiresAuth = true
					break
				}
			}

			// Retrieve the session token from the context
			token, ok := r.Context().Value(SessionIdContextKey).(string)
			if !ok {
				// Token not present in context
				if requiresAuth {
					http.Redirect(w, r, "/signup", http.StatusFound)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			// Check if user ID is already in the context
			if userID, ok := r.Context().Value(UserIdContextKey).(string); ok {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserIdContextKey, userID)))
				return
			}

			// Check cache first
			if userID, found := cache.Get(token); found {
				// Add user ID to the context and proceed
				ctx := context.WithValue(r.Context(), UserIdContextKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Query the database if not in the cache
			var userID string
			log.Println("searching in db, no cache")
			err := db.QueryRow("SELECT userId FROM Sessions WHERE sessionId = ?", token).Scan(&userID)
			if err != nil {
				if err == sql.ErrNoRows {
					if requiresAuth {
						http.Redirect(w, r, "/signup", http.StatusFound)
					} else {
						http.Error(w, "Unauthorized: Invalid session token", http.StatusUnauthorized)
					}
					return
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Database error querying session: %v", err)
				return
			}

			// Cache the session token and user ID
			cache.Set(token, userID)

			// Add user ID to the context and proceed
			ctx := context.WithValue(r.Context(), UserIdContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

/*
	Oui, un client peut avoir besoin d'accéder simultanément à de nombreuses données mises en cache. C'est précisément pourquoi l'utilisation d'une structure concurrente comme `sync.Map` ou d'un autre mécanisme de cache est utile. Voici les raisons pour lesquelles un cache efficace est essentiel dans ce contexte :

	### Besoins du Client pour un Cache Simultané :

1. **Accès Concurentiel:**
   - Plusieurs goroutines (ou requêtes) du client peuvent tenter d'accéder ou de mettre à jour différentes données en même temps.
   - Le cache doit donc être capable de gérer ces accès concurrents sans créer de conflits ou ralentir les performances.

2. **Récupération Rapide:**
   - Le cache doit être capable de fournir rapidement les données nécessaires pour réduire la latence des requêtes.
   - C'est essentiel pour des tâches comme l'authentification, les autorisations, ou le chargement de configurations.

3. **Réduction de la Charge sur la Base de Données:**
   - Un cache bien géré réduit considérablement le nombre d'appels directs à la base de données, ce qui améliore les performances globales.

### Stratégies pour un Cache Simultané :

1. **`sync.Map`:**
   - Pratique pour les scénarios simples où des paires clé-valeur peuvent être stockées en parallèle.

2. **Mutex et RWMutex:**
   - Pour des structures plus complexes (par exemple, des maps imbriquées), des verrous comme `sync.Mutex` ou `sync.RWMutex` peuvent être utilisés pour assurer une synchronisation fine.

3. **Caches Tiers:**
   - Des solutions comme Redis ou Memcached permettent d'avoir un cache distribué accessible depuis plusieurs services ou instances.

### Exemples de Scénarios Simultanés :

1. **Session et Authentification:**
   - Vérification des sessions utilisateur lors de multiples connexions simultanées.

2. **Profil Utilisateur:**
   - Chargement de profils utilisateur et de préférences pour les différents clients.

3. **Données de Configuration:**
   - Fourniture rapide des configurations partagées entre différentes requêtes.

### Résumé :

- Les clients ont souvent besoin d'accéder simultanément à diverses informations.
- Un cache bien conçu peut gérer ces demandes en fournissant un accès rapide et concurrentiel.
- Assurez-vous que votre cache peut répondre aux besoins spécifiques du client tout en restant simple et performant.
*/
