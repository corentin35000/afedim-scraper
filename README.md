# Agency - Scraper

## TODO
- Voir pourquoi on'ai obligé de supprimer le dossier `tmp` parfois pour que ça prends en compte le nouveau code

## 🛠 Tech Stack

- Go (Language)
- Colly (Library)
- CI / CD (Github Actions)
- DockerCompose (Development-Local)
- Kubernetes (Development-Remote, Staging and Production)

<br /><br /><br /><br />

## 📦 Versionning

On utilise la convention SemVer : https://semver.org/lang/fr/ <br /><br />
Pour une Release classique : MAJOR.MINOR.PATCH <br />
Pour une Pre-Release, exemples : MAJOR.MINOR.PATCH-rc.0 OR MAJOR.MINOR.PATCH-beta.3 <br /><br />

Nous utilison release-please de Google pour versionner, via Github Actions. <br />
Pour que cela sois pris en compte il faut utiliser les conventionnal commits : https://www.conventionalcommits.org/en/v1.0.0/ <br />
Release Please crée une demande d'extraction de version après avoir remarqué que la branche par défaut contient des « unités publiables » depuis la dernière version. Une unité publiable est un commit sur la branche avec l'un des préfixes suivants : `feat` / `feat!` et `fix` / `fix!`. <br /><br />

La première Release que créer release-please automatiquement est la version : 1.0.0 <br />
Pour créer une Pre-Release faire un commit vide, par exemple si on'ai à la version 1.0.0, on peut faire :

```bash
git commit --allow-empty -m "chore: release 1.1.0-rc.0" -m "Release-As: 1.1.0-rc.0"
```

<br /><br /><br /><br />

## 🚀 Conventions de Commit

Nous utilisons les conventions de commit pour maintenir une cohérence dans l'historique du code et faciliter le versionnement automatique avec release-please. Voici les types de commits que nous utilisons, ainsi que leur impact sur le versionnage :

- feat : Introduction d'une nouvelle fonctionnalité pour l'utilisateur. Entraîne une augmentation de la version mineure (par exemple, de 1.0.0 à 1.1.0).

- feat! : Introduction d'une nouvelle fonctionnalité avec des modifications incompatibles avec les versions antérieures (breaking changes). Entraîne une augmentation de la version majeure (par exemple, de 1.0.0 à 2.0.0).

- fix : Correction d'un bug pour l'utilisateur. Entraîne une augmentation de la version patch (par exemple, de 1.0.0 à 1.0.1).

- fix! : Correction d'un bug avec des modifications incompatibles avec les versions antérieures (breaking changes). Entraîne une augmentation de la version majeure.

- docs : Changements concernant uniquement la documentation. N'affecte pas la version.

- style : Changements qui n'affectent pas le sens du code (espaces blancs, mise en forme, etc.). N'affecte pas la version.

- refactor : Modifications du code qui n'apportent ni nouvelle fonctionnalité ni correction de bug. N'affecte pas la version.

- perf : Changements de code qui améliorent les performances. Peut entraîner une augmentation de la version mineure.

- test : Ajout ou correction de tests. N'affecte pas la version.

- chore : Changements qui ne modifient ni les fichiers source ni les tests (par exemple, mise à jour des dépendances). N'affecte pas la version.

- ci : Changements dans les fichiers de configuration et les scripts d'intégration continue (par exemple, GitHub Actions). N'affecte pas la version.

- build : Changements qui affectent le système de build ou les dépendances externes (par exemple, npm, Docker). N'affecte pas la version.

- revert : Annulation d'un commit précédent. N'affecte pas la version.

Pour indiquer qu'un commit introduit des modifications incompatibles avec les versions antérieures (breaking changes), ajoutez un ! après le type de commit, par exemple feat! ou fix!.

Pour plus de détails sur les conventions de commit, consultez : [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)

<br /><br /><br /><br />

## 📚 Domains of different environments

- Production : https://test.crzcommon.com
- Staging : https://staging.test.crzcommon.com
- Development-Remote : https://dev.test.crzcommon.com

<br /><br /><br /><br />

## ⚙️ Setup Environment Development

1. Clone the project repository using the following commands :

```bash
git clone git@github.com:corentin35000/crypto-viz-scraper.git
```

2. Steps by Platform :

```bash
# Windows :
1. Requirements : Windows >= 10
2. Download and Install WSL2 : https://learn.microsoft.com/fr-fr/windows/wsl/install
3. Download and Install Docker Desktop : https://docs.docker.com/desktop/install/windows-install/

# macOS :
1. Requirements : macOS Intel x86_64 or macOS Apple Silicon arm64
2. Requirements (2) : macOS 11.0 (Big Sur)
2. Download and Install Docker Desktop : https://docs.docker.com/desktop/install/mac-install/

# Linux (Ubuntu / Debian) :
1. Requirements : Ubuntu >= 20.04 or Debian >= 10
2. Download and Install Docker (Ubuntu) : https://docs.docker.com/engine/install/ubuntu/
3. Download and Install Docker (Debian) : https://docs.docker.com/engine/install/debian/
```

<br /><br /><br /><br />

## 🔄 Cycle Development

1. macOS / Windows : Open Docker Desktop
2. Run command :
```bash
   docker-compose up
```

<br /><br /><br /><br />

## 🚀 Production

### ⚙️➡️ Automatic Distribution Process (CI / CD)

#### Si c'est un nouveau projet suivez les instructions :

1. Ajoutées les SECRETS_GITHUB pour :
   - DOCKER_HUB_USERNAME
   - DOCKER_HUB_ACCESS_TOKEN
   - KUBECONFIG
   - PAT (crée un nouveau token si besoin sur le site de github puis dans le menu du "Profil" puis -> "Settings" -> "Developper Settings' -> 'Personnal Access Tokens' -> Tokens (classic))
