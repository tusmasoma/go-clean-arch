# Clean Architecture
## Description
Clean Architecture is an architectural concept proposed by Robert C. Martin (commonly known as Uncle Bob) in 2012, with the goal of ensuring independence from databases and frameworks. The diagram below represents its iconic image.

![clean_arc](https://github.com/user-attachments/assets/15d93d0c-3a53-46cd-83ba-f394e35cd1ed)

Clean Architecture should be understood not as a "specific architectural pattern," but rather as "principles to follow in architectural design."

Before the introduction of Clean Architecture, approaches like Layered Architecture, Hexagonal Architecture, and Onion Architecture were widely known, all aiming to achieve "separation of concerns." Clean Architecture can be seen as a standardization that integrates these approaches.

As shown in the diagram above, Clean Architecture does not mandate a strict division into four layers: **Enterprise Business Rules**, **Application Business Rules**, **Interface Adapters**, and **Frameworks & Drivers**. This is merely an example. Robert C. Martin himself mentioned that "this diagram is just an overview, and Clean Architecture can have more than four layers."

The crucial point is "separation of concerns." To achieve "separation of concerns," important concepts such as "dependency rules" between layers and "dependency injection" play a vital role.

**Separation of Concerns**
"Separation of concerns" refers to the idea that different elements of a system should have different roles and responsibilities. For example, business logic should exist independently of the database or UI. This enables each part of the system to be modified independently, such as changing the UI without affecting the business logic or changing the database without requiring a complete system rebuild.

**Dependency Rules**
In the "dependency rules," the direction of dependencies in the source code should always point inward (toward the business rules). The inner layers should be designed so that they do not depend on the outer layers (UI, database, frameworks).

The inner layers should know nothing about the outer layers and should avoid referencing specific classes, functions, or variables defined in the outer layers. To achieve this, designs that reverse dependencies using interfaces or abstract classes are crucial.

Moreover, it's important not to depend on data formats or frameworks used in the outer layers. The inner layers should not be influenced by the technologies or formats used externally. This ensures that even if external elements change, the core of the system remains unaffected.

**Dependency Injection**
Dependency Injection is a technique where the dependencies required by a class or component are provided from outside. This prevents a class from creating its dependencies or becoming tied to specific implementations. Using Dependency Injection makes it easier to manage dependencies and keeps the coupling between modules low.

---

By achieving "separation of concerns," the following benefits can be realized:
1. **Framework Independence**: The architecture does not depend on specific frameworks, allowing them to be used as tools. This prevents the system from being constrained by the limitations of frameworks.
2. **Testability**: Business rules can be tested independently of external factors such as UI, database, and web servers.
3. **UI Independence**: The UI can be changed easily without affecting other parts of the system. For example, switching from a web UI to a console UI would not impact the business logic.
4. **Database Independence**: You can easily switch from Oracle or SQL Server to MongoDB, BigTable, CouchDB, etc. The business rules are not dependent on the database.
5. **External Function Independence**: The business rules are independent of external elements and do not need to know anything about the external systems.

## This Project
In this project, we adopt the following four primary layers:
- **Entity Layer**
- **Repository Layer**
- **Usecase Layer**
- **Interfaces Layer**

Here is a brief explanation of each layer.

--- 
### Entity Layer
This layer defines the most critical business rules and domain within the system. From a DDD (Domain-Driven Design) perspective, the Entity Layer is where domain models and important operations on the domain are defined. These entities are independent of other parts of the system and represent universal rules that are not dependent on specific applications or infrastructure.

The Entity Layer corresponds to **Enterprise Business Rules** in the Clean Architecture diagram.

### Repository Layer
This layer is responsible for implementing the technical details, such as databases and email sending. It abstracts the handling of data persistence and interaction with external systems so that other layers do not depend on technical details. For example, the details of database access are hidden within the Repository Layer, and other layers are unaffected by its specific implementation.

The Repository Layer corresponds to **Frameworks & Drivers** in the Clean Architecture diagram.

### Usecase Layer
This layer is responsible for implementing specific use cases (business logic) provided by the system. It defines the system's behavior by utilizing entities and repositories. Use cases represent scenarios or processes to achieve specific business objectives and operate on entities.

Changes in the Usecase Layer are expected not to impact the entities. Additionally, this layer should remain independent of changes in databases, UI, or common frameworks. This layer is isolated from these concerns.

The Usecase Layer corresponds to **Application Business Rules** in the Clean Architecture diagram.

### Interfaces Layer
This layer is responsible for implementing the user interface (UI) and interaction with external systems. It defines the interfaces with the outside world, receives input from external sources, and delegates processing to the internal use cases or entities. It also returns processing results to the external world.

Changes in the Interfaces Layer are expected not to affect other layers.

The Interfaces Layer corresponds to both **Interface Adapters** and **Frameworks & Drivers** in the Clean Architecture diagram.

--- 

### Dependencies
The dependencies between the layers are as follows:

![clean_arc drawio (2)](https://github.com/user-attachments/assets/d8c70210-868e-4f22-91ff-4eb8793171a5)


- Solid Arrows
    - Solid arrows indicate direct dependencies. These are used when one layer directly depends on another layer or component. This dependency might represent situations like:
- Dashed Arrows
    - Dashed arrows indicate indirect dependencies or abstracted dependencies.
        For instance, the Usecase Layer depends on the interfaces in the Repository Layer but does not depend on the specific implementation.

### Dependency Injection
In progress of writing


## HOW To Run
In progress of writing
